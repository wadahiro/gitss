package importer

import (
	"bytes"
	"log"
	"fmt"
	"strings"
	// "io"
	// "io/ioutil"
	// "os"
	// "path/filepath"

	"github.com/wadahiro/gitss/server/indexer"
	"github.com/wadahiro/gitss/server/repo"

	gitm "github.com/gogits/git-module"
	"time"
)

type GitImporter struct {
	dataDir string
	indexer indexer.Indexer
	debug   bool
}

func NewGitImporter(dataDir string, indexer indexer.Indexer, debugMode bool) *GitImporter {
	return &GitImporter{dataDir: dataDir, indexer: indexer, debug: debugMode}
}

func (g *GitImporter) Run(organization string, projectName string, url string) {
	log.Printf("Clone from %s %s %s\n", organization, projectName, url)

	splitedUrl := strings.Split(url, "/")
	repoName := splitedUrl[len(splitedUrl)-1]
	repoPath := fmt.Sprintf("%s/%s/%s/%s", g.dataDir, organization, projectName, repoName)

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
		log.Println(branch)
		g.CreateBranchIndex(repo, branch)
	}
}

func (g *GitImporter) CreateBranchIndex(repo *repo.GitRepo, branchName string) {
	commitId, _ := repo.GetBranchCommitID(branchName)

	if g.debug {
		log.Printf("Indexing start: @%s %s/%s (%s) %s\n", repo.Organization, repo.Project, repo.Repository, branchName, commitId)
	}
	
	start := time.Now()

	// containBranches, _ := ContainsBranch(repo.Path, commitId)

	// if g.debug {
	// 	log.Println("ContainsBranches", containBranches)
	// }

	// commit, err := repo.GetCommit(commitId)

	commit, _ := repo.GetCommit(commitId)
	tree := commit.Tree()

	iter := tree.Files()

	defer func() {
		iter.Close()
	}()

	for true {
		f, err := iter.Next()
		if err != nil {
			break
		}

		if g.debug {
			log.Printf("100644 blob %s %s %d\n", f.Hash, f.Name, f.Size)
		}

		if f.Size > 1024*1024 { // 1MB
			continue
		}

		blobHash := f.Hash.String()

		blob, err := repo.Blob(f.Hash)
		if err != nil {
			continue
		}

		reader, _ := blob.Reader()

		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		content := buf.String()

		g.CreateFileIndex(repo.Organization, repo.Project, repo.Repository, branchName, f.Name, blobHash, content)
	}
	
	end := time.Now()
	time := (end.Sub(start)).Seconds()
	log.Printf("Indexing Complete! [%d seconds]\n", time)
}

func (g *GitImporter) CreateFileIndex(organization string, project string, repo string, branch string, filePath string, blob string, content string) {
	g.indexer.UpsertFileIndex(organization, project, repo, branch, filePath, blob, content)
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
