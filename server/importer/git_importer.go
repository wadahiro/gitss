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

	// "github.com/pkg/errors"
	"github.com/wadahiro/gitss/server/config"
	"github.com/wadahiro/gitss/server/indexer"
	"github.com/wadahiro/gitss/server/repo"

	gitm "github.com/gogits/git-module"
)

type GitImporter struct {
	config  config.Config
	indexer indexer.Indexer
	debug   bool
}

func NewGitImporter(config config.Config, indexer indexer.Indexer) *GitImporter {
	return &GitImporter{config: config, indexer: indexer, debug: config.Debug}
}

func (g *GitImporter) Run(organization string, projectName string, url string) {
	log.Printf("Clone from %s %s %s\n", organization, projectName, url)

	splitedUrl := strings.Split(url, "/")
	repoName := splitedUrl[len(splitedUrl)-1]
	// Drop ".git" from repoName
	splitedRepoNames := strings.Split(repoName, ".git")
	repoName = splitedRepoNames[0]

	repoPath := fmt.Sprintf("%s/%s/%s/%s.git", g.config.GitDataDir, organization, projectName, repoName)

	if err := gitm.Clone(url, repoPath,
		gitm.CloneRepoOptions{Mirror: true}); err != nil {

		// panic(err)
	}

	log.Println("Fetch...")
	FetchAll(repoPath)
	// git.Pull(repoPath, git.PullRemoteOptions{All: true})

	repo, err := repo.NewGitRepo(organization, projectName, repoName, repoPath, g.debug)
	if err != nil {
		panic(err)
	}

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

	if notFound {
		fmt.Printf("New Indexing start: %s\n", tag)
		g.CreateBranchIndex(repo, branchName, latestCommitId)
	} else {
		fmt.Printf("Update Indexing start: %s\n", tag)
		g.UpdateBranchIndex(repo, branchName, indexedCommitId, latestCommitId)
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

func (g *GitImporter) CreateBranchIndex(repo *repo.GitRepo, branchName string, latestCommitId string) {
	addList, err := repo.GetFileEntries(latestCommitId)
	if err == nil {
		g.handleAddBatch(repo, branchName, latestCommitId, addList)
	}
}

func (g *GitImporter) UpdateBranchIndex(repo *repo.GitRepo, branchName string, fromCommitId string, toCommitId string) {
	addList, delList, err := repo.GetDiffList(fromCommitId, toCommitId)
	if err == nil {
		g.handleAddBatch(repo, branchName, toCommitId, addList)
		g.handleDeleteBatch(repo, branchName, toCommitId, delList)
	}
}

func (g *GitImporter) handleAddBatch(repo *repo.GitRepo, branchName string, commitId string, addList []repo.FileEntry) {
	tag := getLoggingTag(repo, branchName, commitId)

	batch := []indexer.FileIndex{}

	var wg sync.WaitGroup
	queue := make(chan indexer.FileIndex, 10)
	done := make(chan struct{})

	go func() {
		for fileIndex := range queue {
			batch = append(batch, fileIndex)

			// Add index
			batch = g.handleBatch(batch, indexer.ADD, 500)
		}
		close(done)
	}()

	for i := range addList {
		addEntry := addList[i]

		wg.Add(1)
		go func() {
			defer wg.Done()

			// check size
			if !g.checkFileSize(addEntry) {
				if g.debug {
					// fmt.Printf("Skipped indexing for size limit. %s %d > %d\n", addEntry.Path, addEntry.Size, g.config.SizeLimit)
				}
				return
			}

			// check contentType
			if ok, _ := g.checkContentType(repo, addEntry); !ok {
				// fmt.Printf("%s skipped. [%s] %s - %s\n", contentType, tag, addEntry.Blob, addEntry.Path)
				return
			}

			content, err := repo.GetBlobContent(addEntry.Blob)
			if err != nil {
				fmt.Printf("Removed among indexing? Skipped [%s] %s - %s\n[%+v]", tag, addEntry.Blob, addEntry.Path, err)
				return
			}

			fileIndex := indexer.FileIndex{
				Blob:    addEntry.Blob,
				Content: string(content),
				Metadata: []indexer.Metadata{
					indexer.Metadata{
						Organization: repo.Organization,
						Project:      repo.Project,
						Repository:   repo.Repository,
						Ref:          branchName,
						Path:         addEntry.Path,
						Ext:          path.Ext(addEntry.Path),
					},
				},
			}

			queue <- fileIndex
		}()
	}

	wg.Wait()
	close(queue)
	<-done

	// Add index remains
	batch = g.handleBatch(batch, indexer.ADD, -1)
}

func (g *GitImporter) checkFileSize(fileEntry repo.FileEntry) bool {
	if fileEntry.Size > g.config.SizeLimit {
		return false
	}
	return true
}

func (g *GitImporter) checkContentType(repo *repo.GitRepo, fileEntry repo.FileEntry) (bool, string) {
	contentType, err := repo.DetectBlobContentType(fileEntry.Blob)
	if err != nil {
		fmt.Println("Failed to read contentType. " + fileEntry.Path)
		return false, ""
	}

	// @TODO Extract text from binary in the future?
	if strings.HasPrefix(contentType, "text/") || contentType == "application/octet-stream" {
		return true, contentType
	} else {
		return false, contentType
	}
}

func (g *GitImporter) handleDeleteBatch(repo *repo.GitRepo, branchName string, commitId string, delList []repo.FileEntry) {
	batch := []indexer.FileIndex{}

	for i := range delList {
		addEntry := delList[i]

		fileIndex := indexer.FileIndex{
			Blob: addEntry.Blob,
			Metadata: []indexer.Metadata{
				indexer.Metadata{
					Organization: repo.Organization,
					Project:      repo.Project,
					Repository:   repo.Repository,
					Ref:          branchName,
					Path:         addEntry.Path,
					Ext:          path.Ext(addEntry.Path),
				},
			},
		}

		batch = append(batch, fileIndex)

		// Delete index
		batch = g.handleBatch(batch, indexer.DELETE, 500)
	}
	// Delete index remains
	batch = g.handleBatch(batch, indexer.DELETE, -1)
}

func (g *GitImporter) handleBatch(batch []indexer.FileIndex, batchMethod indexer.BatchMethod, batchSize int) []indexer.FileIndex {
	if len(batch) > 0 && (batchSize == -1 || len(batch) >= batchSize) {

		// fmt.Printf("Indexed %d files start\n", len(batch))
		g.indexer.BatchFileIndex(batch, batchMethod)
		fmt.Printf("Indexed %d files end\n", len(batch))

		batch = nil
	}
	return batch
}

func getLoggingTag(repo *repo.GitRepo, branchName string, commitId string) string {
	tag := fmt.Sprintf("@%s %s/%s (%s) %s", repo.Organization, repo.Project, repo.Repository, branchName, commitId)
	return tag
}

func (g *GitImporter) CreateFileIndex(organization string, project string, repo string, branch string, filePath string, blob string, content string) {
	fileIndex := indexer.NewFileIndex(blob, organization, project, repo, branch, filePath, content)
	g.indexer.UpsertFileIndex(fileIndex)
}

func FetchAll(repoPath string) error {
	cmd := gitm.NewCommand("fetch")
	cmd.AddArguments("--all")
	cmd.AddArguments("--prune")

	_, err := cmd.RunInDirTimeout(-1, repoPath)
	return err
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
