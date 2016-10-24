package importer

import (
	"fmt"
	"log"
	"strings"
	"sync"
	// "io"
	// "io/ioutil"
	// "os"
	"path"

	"github.com/wadahiro/gitss/server/indexer"
	"github.com/wadahiro/gitss/server/repo"
	"github.com/wadahiro/gitss/server/util"

	"time"

	gitm "github.com/gogits/git-module"
)

type GitImporter struct {
	dataDir   string
	indexer   indexer.Indexer
	sizeLimit int64
	debug     bool
}

func NewGitImporter(dataDir string, indexer indexer.Indexer, sizeLimit int64, debugMode bool) *GitImporter {
	return &GitImporter{dataDir: dataDir, indexer: indexer, sizeLimit: sizeLimit, debug: debugMode}
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
		g.CreateBranchIndex(repo, branch)
	}
}

func (g *GitImporter) CreateBranchIndex(repo *repo.GitRepo, branchName string) {
	commitId, _ := repo.GetBranchCommitID(branchName)

	tag := fmt.Sprintf("@%s %s/%s (%s) %s", repo.Organization, repo.Project, repo.Repository, branchName, commitId)
	fmt.Printf("Indexing start: %s\n", tag)

	start := time.Now()

	// containBranches, _ := ContainsBranch(repo.Path, commitId)

	// if g.debug {
	// 	fmt.Println("ContainsBranches", containBranches)
	// }

	// commit, err := repo.GetCommit(commitId)

	commit, _ := repo.GetCommit(commitId)
	tree := commit.Tree()
	iter := tree.Files()
	defer iter.Close()

	tasks := util.GenWorkers(4)
	wg := &sync.WaitGroup{}
	batch := []indexer.FileIndex{}

	for true {
		f, err := iter.Next()
		if err != nil {
			break
		}

		if g.debug {
			// fmt.Printf("100644 blob %s %s %d\n", f.Hash, f.Name, f.Size)
		}

		if f.Size > g.sizeLimit {
			if g.debug {
				fmt.Printf("Skipped indexing for size limit. %s %d > %d", f.Name, f.Size, g.sizeLimit)
			}
			continue
		}

		blobHash := f.Hash.String()

		// s := time.Now()
		contentType, err := repo.DetectBlobContentType(blobHash)
		if err != nil {
			fmt.Printf("Not found blob. Removed? %s %s - %s\n", tag, blobHash, f.Name)
			continue
		}
		// e := time.Now()
		// t := (e.Sub(s)).Seconds()
		// fmt.Println("Detect time", t)

		if g.debug {
			// fmt.Println("Detected ContentType:", contentType)
		}

		// @TODO Extract text from binary in the future?
		var content []byte
		if strings.HasPrefix(contentType, "text/") {
			content, err = repo.GetBlobContent(blobHash)
			if err != nil {
				fmt.Printf("Not found blob. Removed? %s %s - %s\n", tag, blobHash, f.Name)
				continue
			}
		} else {
			// fmt.Println("Drop content", contentType, f.Name)
		}
		textContent := string(content)

		// fmt.Println(textContent)

		fileIndex := indexer.FileIndex{Blob: blobHash, Content: textContent, Metadata: []indexer.Metadata{indexer.Metadata{Organization: repo.Organization, Project: repo.Project, Repository: repo.Repository, Ref: branchName, Path: f.Name, Ext: path.Ext(f.Name)}}}

		batch = append(batch, fileIndex)

		// g.CreateFileIndex(repo.Organization, repo.Project, repo.Repository, branchName, f.Name, blobHash, content)
		batchSize := 512
		if len(batch) == batchSize {
			wg.Add(1)
			batchCopy := make([]indexer.FileIndex, batchSize)
			copy(batchCopy, batch)
			tasks <- func() {
				defer wg.Done()
				fmt.Printf("Indexed %d files start\n", batchSize)
				g.indexer.BatchFileIndex(batchCopy)
				batch = nil
				fmt.Printf("Indexed %d files end\n", batchSize)
			}
		}
	}
	if len(batch) > 0 {
		fmt.Printf("Indexed %d files start\n", len(batch))
		g.indexer.BatchFileIndex(batch)
		fmt.Printf("Indexed %d files end\n", len(batch))
	}
	wg.Wait()

	end := time.Now()
	time := (end.Sub(start)).Seconds()
	fmt.Printf("Indexing Complete! [%f seconds]\n", time)
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
