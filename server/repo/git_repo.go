package repo

import (

	// "io"
	// "io/ioutil"
	// "os"
	// "path/filepath"
	"fmt"

	"gopkg.in/src-d/go-git.v3/utils/fs"
	// "strings"
	"github.com/wadahiro/gitss/server/config"
	"github.com/wadahiro/gitss/server/util"

	"bytes"
	"io"
	"log"
	"net/http"

	gitm "github.com/gogits/git-module"
	"gopkg.in/src-d/go-git.v3"
	core "gopkg.in/src-d/go-git.v3/core"
)

type GitRepoReader struct {
	gitDataDir string
	debug      bool
}

func NewGitRepoReader(config config.Config) *GitRepoReader {
	reader := &GitRepoReader{gitDataDir: config.GitDataDir, debug: config.Debug}
	return reader
}

func (r *GitRepoReader) GetGitRepo(organization string, project string, repoName string) *GitRepo {
	repoPath := getRepoPath(r.gitDataDir, organization, project, repoName)

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
	Repository   string
	Path         string
	gitmRepo     *gitm.Repository
	repo         *git.Repository
}

type Source struct {
	Offset  int    `json:"offset"`
	Preview string `json:"preview"`
	Hits    []int  `json:"hits"`
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

func (r *GitRepo) GetBlobContent(blobId string) ([]byte, error) {
	blob, err := r.GetBlob(blobId)
	if err != nil {
		return nil, err
	}

	reader, _ := blob.Reader()
	defer reader.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.Bytes(), nil
}

func (r *GitRepo) DetectBlobContentType(blob *git.Blob) (string, error) {
	reader, _ := blob.Reader()

	defer reader.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer[:n])

	return contentType, nil
}

func (r *GitRepo) FilterBlob(blobId string, filter func(line string) bool, before int, after int) []util.TextPreview {
	blob, _ := r.GetBlob(blobId)
	reader, _ := blob.Reader()
	defer reader.Close()

	previews := util.FilterTextPreview(reader, filter, before, after)

	return previews
}

type FileEntry struct {
	Blob string
	Path string
	Size int64
}

func (r *GitRepo) GetFileEntries(commitId string) ([]FileEntry, error) {
	commit, err := r.GetCommit(commitId)
	if err != nil {
		return nil, err
	}

	tree := commit.Tree()
	iter := tree.Files()
	defer iter.Close()

	list := []FileEntry{}

	for true {
		f, err := iter.Next()
		if err != nil {
			break
		}

		blob := f.Hash.String()
		list = append(list, FileEntry{Blob: blob, Path: f.Name, Size: f.Size})
	}
	return list, nil
}

func (r *GitRepo) GetDiffList(from string, to string) ([]FileEntry, []FileEntry, error) {
	stdout, err := gitm.NewCommand("diff", "--raw", "--abbrev=40", "-z", from, to).RunInDirTimeout(-1, r.Path)
	if err != nil {
		return nil, nil, err
	}

	addList := []FileEntry{}
	delList := []FileEntry{}

	parts := bytes.Split(stdout, []byte{'\u0000'})

	for index := 0; index < len(parts); index++ {
		row := parts[index]
		cols := bytes.Split(row, []byte{' '})

		if len(cols) <= 6 {
			continue
		}

		oldBlob := string(cols[2])
		newBlob := string(cols[2])
		status := string(cols[4])[0:1]

		// See 'Possible status letters' in https://git-scm.com/docs/git-diff
		switch status {
		case "A": // Add case
			index++
			path := string(parts[index])
			addList = append(addList, FileEntry{Blob: newBlob, Path: path})

		case "C": // Copy case
			index = index + 2
			toPath := string(parts[index])
			addList = append(addList, FileEntry{Blob: newBlob, Path: toPath})

		case "D": // Delete case
			index++
			path := string(parts[index])
			delList = append(delList, FileEntry{Blob: oldBlob, Path: path})

		case "M": // Modify case
		case "T": // Change in the type of the file case
			index++
			if oldBlob != newBlob {
				path := string(parts[index])
				delList = append(delList, FileEntry{Blob: oldBlob, Path: path})
				addList = append(addList, FileEntry{Blob: newBlob, Path: path})
			}

		case "R": // Rename case
			index++
			fromPath := string(parts[index])
			index++
			toPath := string(parts[index])
			delList = append(delList, FileEntry{Blob: oldBlob, Path: fromPath})
			addList = append(addList, FileEntry{Blob: newBlob, Path: toPath})

		case "U": // Unmerge case
			continue
		case "X": // Delete case
			log.Println("unknown change type (most probably a git bug, please report it)")
		}
	}

	return addList, delList, nil
}

func getRepoPath(gitDataDir string, organization string, project string, repoName string) string {
	repoPath := fmt.Sprintf("%s/%s/%s/%s.git", gitDataDir, organization, project, repoName)
	return repoPath
}
