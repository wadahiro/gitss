package importer

import (
	"fmt"
	"log"
	"strings"
	// "io"
	// "io/ioutil"
	// "os"
	"path"

	"time"

	"github.com/pkg/errors"
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
	repoPath := fmt.Sprintf("%s/%s/%s/%s", g.config.GitDataDir, organization, projectName, repoName)

	// Drop ".git" from repoName
	splitedRepoNames := strings.Split(repoName, ".git")
	if len(splitedRepoNames) > 1 {
		repoName = splitedRepoNames[0]
	}

	if err := gitm.Clone(url, repoPath,
		gitm.CloneRepoOptions{Mirror: true}); err != nil {

		// panic(err)
	}

	log.Println("Fetch...")
	FetchAll(repoPath)
	// git.Pull(repoPath, git.PullRemoteOptions{All: true})

	repo, err := repo.NewGitRepo(organization, projectName, repoName, repoPath)
	if err != nil {
		panic(err)
	}

	branches, _ := repo.GetBranches()

	for _, branch := range branches {
		g.RunIndexing(repo, branch)
	}
}

func (g *GitImporter) RunIndexing(repo *repo.GitRepo, branchName string) {
	latestCommitId, _ := repo.GetBranchCommitID(branchName)
	indexedCommitId, notFound := g.config.GetIndexedCommitID(config.LatestIndex{
		Organization: repo.Organization,
		Project:      repo.Project,
		Repository:   repo.Repository,
		Ref:          branchName,
	})

	tag := getLoggingTag(repo, branchName, latestCommitId)
	fmt.Printf("Indexing start: %s\n", tag)

	start := time.Now()

	if notFound {
		g.CreateBranchIndex(repo, branchName, latestCommitId)
	} else {
		g.UpdateBranchIndex(repo, branchName, indexedCommitId, latestCommitId)
	}

	// Save config after index completed
	g.config.UpdateLatestIndex(config.LatestIndex{
		Organization: repo.Organization,
		Project:      repo.Project,
		Repository:   repo.Repository,
		Ref:          branchName,
	}, latestCommitId)

	end := time.Now()
	time := (end.Sub(start)).Seconds()
	fmt.Printf("Indexing Complete! [%f seconds]\n", time)
}

func (g *GitImporter) CreateBranchIndex(repo *repo.GitRepo, branchName string, latestCommitId string) {
	addList, err := repo.GetFileEntries(latestCommitId)
	if err != nil {
		g.handleAddBatch(repo, branchName, latestCommitId, addList)
	}
}

func (g *GitImporter) UpdateBranchIndex(repo *repo.GitRepo, branchName string, fromCommitId string, toCommitId string) {
	addList, delList, err := repo.GetDiffList(fromCommitId, toCommitId)
	if err != nil {
		g.handleAddBatch(repo, branchName, toCommitId, addList)
		g.handleDeleteBatch(repo, branchName, toCommitId, delList)
	}
}

func (g *GitImporter) handleAddBatch(repo *repo.GitRepo, branchName string, commitId string, addList []repo.FileEntry) {
	tag := getLoggingTag(repo, branchName, commitId)

	batch := []indexer.FileIndex{}

	for i := range addList {
		addEntry := addList[i]

		textContent, err := g.getFileContent(repo, addEntry.Blob, addEntry.Path)
		if err != nil {
			if g.debug {
				fmt.Printf("Skipped [%s] %s %s - %s\n[%+v]", tag, addEntry.Blob, addEntry.Path, err)
			}
			continue
		}

		fileIndex := indexer.FileIndex{
			Blob:    addEntry.Blob,
			Content: textContent,
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

		// Add index
		batch = g.handleBatch(batch, indexer.ADD, 500)
	}
	// Add index remains
	batch = g.handleBatch(batch, indexer.ADD, -1)
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
	if batchSize > 0 && len(batch) == batchSize {

		fmt.Printf("Indexed %d files start\n", len(batch))
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

func (g *GitImporter) getFileContent(repo *repo.GitRepo, blob string, filePath string) (string, error) {
	b, err := repo.GetBlob(blob)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read blob. "+filePath)
	}

	if b.Size > g.config.SizeLimit {
		return "", errors.Errorf("Skipped indexing for size limit. %s %d > %d", filePath, b.Size, g.config.SizeLimit)
	}

	// s := time.Now()
	contentType, err := repo.DetectBlobContentType(b)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read contentType. "+filePath)
	}

	// e := time.Now()
	// t := (e.Sub(s)).Seconds()
	// fmt.Println("Detect time", t)

	// fmt.Println("Detected ContentType:", contentType)

	// @TODO Extract text from binary in the future?
	var content []byte
	if strings.HasPrefix(contentType, "text/") {
		content, err = repo.GetBlobContent(blob)
		if err != nil {
			return "", errors.Wrap(err, "Failed to read blob. "+filePath)
		}
	} else {
		// fmt.Println("Drop content", contentType, f.Name)
	}
	textContent := string(content)

	return textContent, nil
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
