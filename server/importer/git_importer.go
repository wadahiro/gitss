package importer

import (
	"fmt"
	"log"
	"strings"
	// "io"
	// "io/ioutil"
	// "os"
	"sync"

	"time"

	"github.com/pkg/errors"
	"github.com/wadahiro/gitss/server/config"
	"github.com/wadahiro/gitss/server/indexer"
	"github.com/wadahiro/gitss/server/repo"
	"github.com/wadahiro/gitss/server/util"
)

type GitImporter struct {
	config  *config.Config
	indexer indexer.Indexer
	reader  *repo.GitRepoReader
	debug   bool
}

func NewGitImporter(config *config.Config, indexer indexer.Indexer) *GitImporter {
	r := repo.NewGitRepoReader(config)
	return &GitImporter{config: config, indexer: indexer, reader: r, debug: config.Debug}
}

func (g *GitImporter) Run(organization string, project string, url string) {
	log.Printf("Clone from %s %s %s\n", organization, project, url)

	repo, err := g.reader.CloneGitRepo(organization, project, url)
	if err != nil {
		log.Printf("Not found the repository: %s:%s/%s %+v\n", organization, project, url, err)
		return
	}

	repo.FetchAll()

	log.Printf("Fetched all. %s %s %s \n", organization, project, url)

	// branches in the git repository
	branches, _ := repo.GetBranches()

	// branches in the config file
	refs, ok := g.config.GetRefs(organization, project, repo.Repository)
	if !ok {
		log.Printf("Not found repository setting. %s:%s/%s \n", organization, project, repo.Repository)
		return
	}

	log.Printf("Start indexing for %s:%s/%s %v -> %v\n", organization, project, repo.Repository, refs, branches)

	start := time.Now()

	branchCommitIdMap, err := repo.GetBrancheCommitIdMap()
	if err != nil {
		log.Printf("Not found branches. %s:%s/%s %+v\n", organization, project, repo.Repository, err)
		return
	}

	for _, branch := range branches {
		g.RunIndexing(branchCommitIdMap, url, repo, branch)
	}

	// Remove index for removed branches
	removeRefs := []string{}
	for _, ref := range refs {
		found := false
		for _, branch := range branches {
			if ref.Name == branch {
				found = true
				break
			}
		}
		if !found {
			removeRefs = append(removeRefs, ref.Name)
		}
	}

	if len(removeRefs) > 0 {
		log.Printf("Start index deleting for %s:%s/%s %v\n", organization, project, repo.Repository, removeRefs)
		g.indexer.DeleteIndexByRefs(organization, project, repo.Repository, removeRefs)

		// Save config after deleting index completed
		g.config.DeleteLatestIndexRefs(organization, project, repo.Repository, removeRefs)
	}

	end := time.Now()
	time := (end.Sub(start)).Seconds()
	log.Printf("Indexing Complete! [%f seconds] for %s:%s/%s\n", time, organization, project, repo.Repository)
}

func (g *GitImporter) RunIndexing(branchCommitIdMap map[string]string, url string, repo *repo.GitRepo, branchName string) error {
	latestCommitId, _ := branchCommitIdMap[branchName]
	indexedCommitId, ok := g.config.GetIndexedCommitID(repo.Organization, repo.Project, repo.Repository, branchName)

	tag := getLoggingTag(repo, branchName, latestCommitId)

	if latestCommitId == indexedCommitId {
		log.Printf("Already up-to-date. %s", tag)
		return nil
	}

	queue := make(chan indexer.FileIndexOperation, 100)

	if !ok {
		fmt.Printf("New Indexing start: %s\n", tag)
		g.CreateBranchIndex(branchCommitIdMap, queue, repo, branchName, latestCommitId)
	} else {
		fmt.Printf("Update Indexing start: %s\n", tag)
		g.UpdateBranchIndex(queue, repo, branchName, indexedCommitId, latestCommitId)
	}

	callBach := func(operations []indexer.FileIndexOperation) {
		err := g.indexer.BatchFileIndex(operations)
		if err != nil {
			errors.Errorf("Batch indexed error: %+v", err)
		} else {
			fmt.Printf("Batch indexed %d files.\n", len(operations))
		}
	}

	// batch
	operations := []indexer.FileIndexOperation{}
	var opsSize int64 = 0
	var batchLimitSize int64 = 1024 * 1024 // 1MB

	for op := range queue {
		operations = append(operations, op)
		opsSize += op.FileIndex.Size

		// show progress
		if len(operations)%80 == 0 {
			fmt.Printf("\n")
		}
		fmt.Printf(".")

		if opsSize >= batchLimitSize {
			fmt.Printf("\n")

			callBach(operations)

			// reset
			operations = nil
			opsSize = 0
		}
	}

	// remains
	if len(operations) > 0 {
		fmt.Printf("\n")
		callBach(operations)
	}

	// Save config after index completed
	g.config.UpdateLatestIndex(url, repo.Organization, repo.Project, repo.Repository, branchName, latestCommitId)

	return nil
}

func (g *GitImporter) CreateBranchIndex(branchCommitIdMap map[string]string, queue chan indexer.FileIndexOperation, r *repo.GitRepo, branchName string, latestCommitId string) error {

	workers := util.GenWorkers(10)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := r.GetFileEntriesIterator(latestCommitId, func(fileEntry repo.FileEntry) {
			// check size
			if fileEntry.Size > g.config.SizeLimit {
				return
			}
			wg.Add(1)
			workers <- func() {
				defer wg.Done()
				// heavy process

				fileIndex := indexer.FileIndex{
					Metadata: indexer.Metadata{
						Blob:         fileEntry.Blob,
						Organization: r.Organization,
						Project:      r.Project,
						Repository:   r.Repository,
						Path:         fileEntry.Path,
						Ext:          indexer.GetExt(fileEntry.Path),
						Size:         fileEntry.Size,
					},
				}

				if g.config.PreFetchRefs {
					ok, _ := g.indexer.Exists(fileIndex)
					if ok {
						log.Println("Already indexed by preFetchRefs.")
						return
					}
				}

				// check contentType and retrive the file content
				// !! this will be heavy process !!
				ok, content := g.checkContentType(r, fileEntry)
				if !ok {
					// fmt.Printf("%s skipped. [%s] %s - %s\n", contentType, tag, addEntry.Blob, addEntry.Path)
					return
				}

				// check whether the same file exists in the other branches
				refs := []string{branchName}
				if g.config.PreFetchRefs {
					for otherBranch, commitId := range branchCommitIdMap {
						if otherBranch != branchName {
							exists, _ := r.ExistsInCommit(commitId, fileEntry.Path, fileEntry.Blob)
							if exists {
								refs = append(refs, otherBranch)
							}
						}
					}
					log.Println("PreFetchRefs: ", refs)
				}

				fileIndex.Refs = refs
				fileIndex.Content = content

				queue <- indexer.FileIndexOperation{Method: indexer.ADD, FileIndex: fileIndex}
			}
		})
		if err != nil {
			log.Printf("NotFound commitId: %s, %+v", latestCommitId, err)
		}
	}()

	go func() {
		wg.Wait()
		close(queue)
	}()

	return nil
}

func (g *GitImporter) UpdateBranchIndex(queue chan indexer.FileIndexOperation, r *repo.GitRepo, branchName string, fromCommitId string, toCommitId string) error {

	workers := util.GenWorkers(10)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := r.GetDiffEntriesIterator(fromCommitId, toCommitId, func(fileEntry repo.FileEntry, status string) {
			wg.Add(1)
			workers <- func() {
				defer wg.Done()
				// heavy process

				if status == "A" {
					// check size
					size, err := r.GetBlobSize(fileEntry.Blob)
					if err != nil {
						fmt.Println("Failed to read size. " + fileEntry.Path)
						return
					}
					if size > g.config.SizeLimit {
						return
					}
					fileEntry.Size = size

					// check contentType
					ok, content := g.checkContentType(r, fileEntry)
					if !ok {
						// fmt.Printf("%s skipped. [%s] %s - %s\n", contentType, tag, addEntry.Blob, addEntry.Path)
						return
					}

					fileIndex := indexer.FileIndex{
						Metadata: indexer.Metadata{
							Blob:         fileEntry.Blob,
							Organization: r.Organization,
							Project:      r.Project,
							Repository:   r.Repository,
							Refs:         []string{branchName},
							Path:         fileEntry.Path,
							Ext:          indexer.GetExt(fileEntry.Path),
							Size:         fileEntry.Size,
						},
						Content: content,
					}
					// Add index
					queue <- indexer.FileIndexOperation{Method: indexer.ADD, FileIndex: fileIndex}

				} else {
					fileIndex := indexer.FileIndex{
						Metadata: indexer.Metadata{
							Blob:         fileEntry.Blob,
							Organization: r.Organization,
							Project:      r.Project,
							Repository:   r.Repository,
							Refs:         []string{branchName},
							Path:         fileEntry.Path,
							Ext:          indexer.GetExt(fileEntry.Path),
							Size:         fileEntry.Size,
						},
					}
					// Delete index
					queue <- indexer.FileIndexOperation{Method: indexer.DELETE, FileIndex: fileIndex}
				}
			}
		})
		if err != nil {
			log.Printf("NotFound diff: %s..%s %+v", fromCommitId, toCommitId, err)
		}
	}()

	go func() {
		wg.Wait()
		close(queue)
	}()

	return nil
}

func (g *GitImporter) checkContentType(repo *repo.GitRepo, fileEntry repo.FileEntry) (bool, string) {
	contentType, content, err := repo.DetectBlobContentType(fileEntry.Blob)
	if err != nil {
		fmt.Println("Failed to read contentType. " + fileEntry.Path)
		return false, ""
	}

	// @TODO Extract text from binary in the future?
	if strings.HasPrefix(contentType, "text/") || contentType == "application/octet-stream" {
		return true, string(content)
	} else {
		return false, ""
	}
}

func getLoggingTag(repo *repo.GitRepo, branchName string, commitId string) string {
	tag := fmt.Sprintf("%s:%s/%s (%s) [%s]", repo.Organization, repo.Project, repo.Repository, branchName, commitId)
	return tag
}
