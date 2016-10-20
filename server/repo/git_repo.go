package repo

import (

	// "io"
	// "io/ioutil"
	// "os"
	// "path/filepath"
	"fmt"
	"gopkg.in/src-d/go-git.v3/utils/fs"
	// "strings"
	"github.com/wadahiro/gitss/server/util"

	gitm "github.com/gogits/git-module"
	"gopkg.in/src-d/go-git.v3"
	core "gopkg.in/src-d/go-git.v3/core"
	"log"
)

type GitRepoReader struct {
	dataDir   string
	debugMode bool
}

func NewGitRepoReader(dataDir string, debugMode bool) *GitRepoReader {
	reader := &GitRepoReader{dataDir: dataDir, debugMode: debugMode}
	return reader
}

func (r *GitRepoReader) GetGitRepo(organization string, project string, repoName string) *GitRepo {

	repoPath := fmt.Sprintf("%s/%s/%s/%s.git", r.dataDir, organization, project, repoName)

	repo, err := NewGitRepo(organization, project, repoName, repoPath)
	if err != nil {
		log.Println(err, repoPath)
		panic(err)
	}
	return repo
}

type GitRepo struct {
	Organization string
	Project      string
	Repository         string
	Path         string
	gitmRepo     *gitm.Repository
	repo         *git.Repository
}

func NewGitRepo(organization string, projectName string, repoName string, repoPath string) (*GitRepo, error) {
	gitmRepo, err := gitm.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}

	fs := fs.NewOS()
	r, _ := git.NewRepositoryFromFS(fs, repoPath)

	return &GitRepo{Organization: organization, Project: projectName, Repository: repoName, Path: repoPath, gitmRepo: gitmRepo, repo: r}, nil
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

func (r *GitRepo) FilterBlob(blobId string, filter func(line string) bool, before int, after int) []util.TextPreview {
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
