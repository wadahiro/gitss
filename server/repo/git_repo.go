package repo

import (

	// "io"
	// "io/ioutil"
	// "os"
	// "path/filepath"
	"fmt"
	// "strings"
	"github.com/wadahiro/GitS/server/util"

	gitm "github.com/gogits/git-module"
	"gopkg.in/src-d/go-git.v4"
	core "gopkg.in/src-d/go-git.v4/core"
)

type GitRepoReader struct {
	dataDir   string
	debugMode bool
}

func NewGitRepoReader(dataDir string, debugMode bool) *GitRepoReader {
	reader := &GitRepoReader{dataDir: dataDir, debugMode: debugMode}
	return reader
}

func (r *GitRepoReader) GetGitRepo(project string, repoName string) *GitRepo {

	repoPath := fmt.Sprintf("%s/%s/%s.git", r.dataDir, project, repoName)

	repo, err := NewGitRepo(project, project, repoPath)
	if err != nil {
		panic(err)
	}
	return repo
}

type GitRepo struct {
	Project  string
	Repo     string
	Path     string
	gitmRepo *gitm.Repository
	repo     *git.Repository
}

func NewGitRepo(projectName string, repoName string, repoPath string) (*GitRepo, error) {
	gitmRepo, err := gitm.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}

	repo, _ := git.NewFilesystemRepository(repoPath)

	return &GitRepo{Project: projectName, Repo: repoName, Path: repoPath, gitmRepo: gitmRepo, repo: repo}, nil
}

func (r *GitRepo) GetBranches() ([]string, error) {
	return r.gitmRepo.GetBranches()
}

func (r *GitRepo) GetBranchCommitID(name string) (string, error) {
	return r.gitmRepo.GetBranchCommitID(name)
}

func (r *GitRepo) GetCommit(commitId string) (*git.Commit, error) {
	return r.repo.Commit(core.NewHash(commitId))
}

func (r *GitRepo) GetBlob(blobId string) (*git.Blob, error) {
	return r.repo.Blob(core.NewHash(blobId))
}

func (r *GitRepo) Blob(hash core.Hash) (*git.Blob, error) {
	return r.repo.Blob(hash)
}

func (r *GitRepo) FilterBlob(blobId string, filter func(line string) bool, before int, after int) []*util.TextPreview {
	blob, _ := r.GetBlob(blobId)
	reader, _ := blob.Reader()

	previews := util.FilterTextPreview(reader, filter, before, after)

	return previews
}

type Source struct {
	Offset  int    `json:"offset"`
	Preview string `json:"preview"`
	Hits    []int  `json:"hits"`
}
