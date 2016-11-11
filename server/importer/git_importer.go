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
	// "github.com/wadahiro/gitss/server/util"
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

	// branches and tags in the git repository
	// @TODO include, exclude options
	branchMap, tagMap, err := repo.GetLatestCommitIdsMap(nil, nil, nil, nil)
	if err != nil {
		log.Printf("Failed to get latest commitIds of branch and tag. %+v\n", err)
		return
	}

	// branches in the config file
	indexed := g.config.GetIndexed(organization, project, repo.Repository)

	log.Printf("Start indexing for %s:%s/%s branches: %v -> %v, tags: %v -> %v\n", organization, project, repo.Repository, indexed.Branches, branchMap, indexed.Tags, tagMap)

	start := time.Now()

	err = g.runIndexing(repo, url, indexed, branchMap, tagMap)
	if err != nil {
		log.Printf("Failed to index. %+v", err)
		return
	}

	// Remove index for removed branches
	removeBranches := []string{}
	for ref, _ := range indexed.Branches {
		found := false
		for branch := range branchMap {
			if ref == branch {
				found = true
				break
			}
		}
		if !found {
			removeBranches = append(removeBranches, ref)
		}
	}

	// Remove index for removed tags
	removeTags := []string{}
	for ref, _ := range indexed.Tags {
		found := false
		for branch := range branchMap {
			if ref == branch {
				found = true
				break
			}
		}
		if !found {
			removeTags = append(removeTags, ref)
		}
	}

	if len(removeBranches) > 0 || len(removeTags) > 0 {
		log.Printf("Start index deleting for %s:%s/%s (%v) (%v)\n", organization, project, repo.Repository, removeBranches, removeTags)
		g.indexer.DeleteIndexByRefs(organization, project, repo.Repository, removeBranches, removeTags)

		// Save config after deleting index completed
		g.config.DeleteIndexed(organization, project, repo.Repository, removeBranches, removeTags)
	}

	end := time.Now()
	time := (end.Sub(start)).Seconds()
	log.Printf("Indexing Complete! [%f seconds] for %s:%s/%s\n", time, organization, project, repo.Repository)
}

func (g *GitImporter) runIndexing(repo *repo.GitRepo, url string, indexed config.Indexed, branchMap config.BrancheIndexedMap, tagMap config.TagIndexedMap) error {
	// collect create file entries
	createBranches := make(map[string]string)
	updateBranches := make(map[string][2]string)
	for branch, latestCommitID := range branchMap {
		found := false
		for indexedBranch, prevCommitID := range indexed.Branches {
			if branch == indexedBranch {
				found = true
				if latestCommitID == prevCommitID {
					log.Printf("Already up-to-date. %s", getLoggingTag(repo, branch, latestCommitID))
				} else {
					updateBranches[branch] = [2]string{prevCommitID, latestCommitID}
				}
				break
			}
		}
		if !found {
			createBranches[branch] = latestCommitID
		}
	}

	createTags := make(map[string]string)
	updateTags := make(map[string][2]string)
	for tag, latestCommitID := range tagMap {
		found := false
		for indexedTag, prevCommitID := range indexed.Tags {
			if tag == indexedTag {
				found = true
				if latestCommitID == prevCommitID {
					log.Printf("Already up-to-date. %s", getLoggingTag(repo, tag, latestCommitID))
				} else {
					updateTags[tag] = [2]string{prevCommitID, latestCommitID}
				}
				break
			}
		}
		if !found {
			createTags[tag] = latestCommitID
		}
	}

	queue := make(chan indexer.FileIndexOperation, 100)

	// process
	g.UpsertIndex(queue, repo, createBranches, createTags, updateBranches, updateTags)

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
	err := g.config.UpdateIndexed(config.Indexed{Organization: repo.Organization, Project: repo.Project, Repository: repo.Repository, Branches: branchMap, Tags: tagMap})

	if err != nil {
		return errors.Wrapf(err, "Faild to update indexed.")
	}

	return nil
}

func (g *GitImporter) UpsertIndex(queue chan indexer.FileIndexOperation, r *repo.GitRepo, branchMap map[string]string, tagMap map[string]string, updateBranchMap map[string][2]string, updateTagMap map[string][2]string) error {
	addFiles, err := r.GetFileEntriesMap(branchMap, tagMap)
	if err != nil {
		return errors.Wrapf(err, "Failed to get file entries. branches: %v tags: %v", branchMap, tagMap)
	}

	updateAddFiles, delFiles, err := r.GetDiffFileEntriesMap(updateBranchMap, updateTagMap)
	if err != nil {
		return errors.Wrapf(err, "Failed to get diff. branches: %v tags: %v", branchMap, tagMap)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		g.handleAddFiles(queue, r, addFiles)
		g.handleAddFiles(queue, r, updateAddFiles)
		g.handleDelFiles(queue, r, delFiles)
	}()

	go func() {
		wg.Wait()
		close(queue)
	}()

	return nil
}

func (g *GitImporter) handleAddFiles(queue chan indexer.FileIndexOperation, r *repo.GitRepo, addFiles map[string]repo.GitFile) {
	for blob, file := range addFiles {
		// check size
		if file.Size > g.config.SizeLimit {
			continue
		}
		for path, loc := range file.Locations {
			// check contentType and retrive the file content
			// !! this will be heavy process !!
			contentType, content, err := g.parseContent(r, blob)
			if err != nil {
				log.Printf("Failed to parse file. [%s] - %s %+v\n", blob, path, err)
				continue
				// return errors.Wrapf(err, "Failed to parse file. [%s] - %s\n", blob, path)
			}

			// @TODO Extract text from binary in the future?
			if !strings.HasPrefix(contentType, "text/") && contentType != "application/octet-stream" {
				continue
			}

			fileIndex := indexer.FileIndex{
				Metadata: indexer.Metadata{
					Blob:         blob,
					Organization: r.Organization,
					Project:      r.Project,
					Repository:   r.Repository,
					Branches:     loc.Branches,
					Tags:         loc.Tags,
					Path:         path,
					Ext:          indexer.GetExt(path),
					Size:         file.Size,
				},
				Content: content,
			}
			queue <- indexer.FileIndexOperation{Method: indexer.ADD, FileIndex: fileIndex}
		}
	}
}

func (g *GitImporter) handleDelFiles(queue chan indexer.FileIndexOperation, r *repo.GitRepo, delFiles map[string]repo.GitFile) {
	for blob, file := range delFiles {
		for path, loc := range file.Locations {
			fileIndex := indexer.FileIndex{
				Metadata: indexer.Metadata{
					Blob:         blob,
					Organization: r.Organization,
					Project:      r.Project,
					Repository:   r.Repository,
					Branches:     loc.Branches,
					Tags:         loc.Tags,
					Path:         path,
				},
			}
			// Delete index
			queue <- indexer.FileIndexOperation{Method: indexer.DELETE, FileIndex: fileIndex}
		}
	}
}

func (g *GitImporter) parseContent(repo *repo.GitRepo, blob string) (string, string, error) {
	contentType, content, err := repo.DetectBlobContentType(blob)
	if err != nil {
		return "", "", errors.Wrapf(err, "Failed to read contentType. %s", blob)
	}
	return string(contentType), string(content), nil
}

func getLoggingTag(repo *repo.GitRepo, ref string, commitId string) string {
	tag := fmt.Sprintf("%s:%s/%s (%s) [%s]", repo.Organization, repo.Project, repo.Repository, ref, commitId)
	return tag
}
