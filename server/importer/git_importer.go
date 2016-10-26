package importer

import (
	"fmt"
	"log"
	"strings"
	// "io"
	// "io/ioutil"
	// "os"
	"path"
	"sync"

	"time"

	"github.com/pkg/errors"
	"github.com/wadahiro/gitss/server/config"
	"github.com/wadahiro/gitss/server/indexer"
	"github.com/wadahiro/gitss/server/repo"
	"github.com/wadahiro/gitss/server/util"

	gitm "github.com/gogits/git-module"
)

type GitImporter struct {
	config  config.Config
	indexer indexer.Indexer
	reader  *repo.GitRepoReader
	debug   bool
}

func NewGitImporter(config config.Config, indexer indexer.Indexer) *GitImporter {
	r := repo.NewGitRepoReader(config)
	return &GitImporter{config: config, indexer: indexer, reader: r, debug: config.Debug}
}

func (g *GitImporter) Run(organization string, projectName string, url string) {
	log.Printf("Clone from %s %s %s\n", organization, projectName, url)

	repo, err := g.reader.CloneGitRepo(organization, projectName, url)
	if err != nil {
		log.Printf("Not found the repository: %s %s %s %+v\n", organization, projectName, url, err)
		return
	}

	repo.FetchAll()

	log.Println("Fetched all.")

	branches, _ := repo.GetBranches()

	for _, branch := range branches {
		g.RunIndexing(url, repo, branch)
	}
}

func (g *GitImporter) RunIndexing(url string, repo *repo.GitRepo, branchName string) {
	latestCommitId, _ := repo.GetBranchCommitID(branchName)
	indexedCommitId, notFound := g.config.GetIndexedCommitID(config.LatestIndex{
		Organization: repo.Organization,
		Project:      repo.Project,
		Repository:   repo.Repository,
		Ref:          branchName,
	})

	tag := getLoggingTag(repo, branchName, latestCommitId)

	start := time.Now()

	queue := make(chan indexer.FileIndexOperation, 100)

	if notFound {
		fmt.Printf("New Indexing start: %s\n", tag)
		g.CreateBranchIndex(queue, repo, branchName, latestCommitId)
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
	for op := range queue {
		operations = append(operations, op)
		opsSize := len(operations)

		// show progress
		if opsSize%100 == 0 {
			fmt.Printf("\n")
		}
		fmt.Printf(".")

		if opsSize >= 1000 {
			fmt.Printf("\n")

			callBach(operations)
			operations = nil
		}
	}

	// remains
	opsSize := len(operations)
	if opsSize > 0 {
		fmt.Printf("\n")
		callBach(operations)
	}

	// Save config after index completed
	g.config.UpdateLatestIndex(url, config.LatestIndex{
		Organization: repo.Organization,
		Project:      repo.Project,
		Repository:   repo.Repository,
		Ref:          branchName,
	}, latestCommitId)

	end := time.Now()
	time := (end.Sub(start)).Seconds()
	fmt.Printf("Indexing Complete! %s [%f seconds]\n", tag, time)
}

func (g *GitImporter) CreateBranchIndex(queue chan indexer.FileIndexOperation, r *repo.GitRepo, branchName string, latestCommitId string) error {

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

				// check contentType
				ok, content := g.checkContentType(r, fileEntry)
				if !ok {
					// fmt.Printf("%s skipped. [%s] %s - %s\n", contentType, tag, addEntry.Blob, addEntry.Path)
					return
				}

				fileIndex := indexer.FileIndex{
					Blob:    fileEntry.Blob,
					Content: content,
					Metadata: []indexer.Metadata{
						indexer.Metadata{
							Organization: r.Organization,
							Project:      r.Project,
							Repository:   r.Repository,
							Ref:          branchName,
							Path:         fileEntry.Path,
							Ext:          path.Ext(fileEntry.Path),
						},
					},
				}
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
						Blob:    fileEntry.Blob,
						Content: content,
						Metadata: []indexer.Metadata{
							indexer.Metadata{
								Organization: r.Organization,
								Project:      r.Project,
								Repository:   r.Repository,
								Ref:          branchName,
								Path:         fileEntry.Path,
								Ext:          path.Ext(fileEntry.Path),
							},
						},
					}
					// Add index
					queue <- indexer.FileIndexOperation{Method: indexer.ADD, FileIndex: fileIndex}

				} else {
					fileIndex := indexer.FileIndex{
						Blob: fileEntry.Blob,
						Metadata: []indexer.Metadata{
							indexer.Metadata{
								Organization: r.Organization,
								Project:      r.Project,
								Repository:   r.Repository,
								Ref:          branchName,
								Path:         fileEntry.Path,
								Ext:          path.Ext(fileEntry.Path),
							},
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
	tag := fmt.Sprintf("@%s %s/%s (%s) %s", repo.Organization, repo.Project, repo.Repository, branchName, commitId)
	return tag
}

func ContainsBranch(repoPath string, commitId string) ([]string, error) {
	cmd := gitm.NewCommand("branch")
	cmd.AddArguments("--contains", commitId)

	stdout, err := cmd.RunInDir(repoPath)

	// log.Println("--------------->", err)
	if err != nil {
		return nil, err
	}
	// log.Println("--------------->", stdout)

	infos := strings.Split(stdout, "\n")
	// log.Println(len(infos))
	branches := make([]string, len(infos)-1)
	for i, info := range infos[:len(infos)-1] {
		branches[i] = info[2:]
	}
	return branches, nil
}
