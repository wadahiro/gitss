package repo

import (

	// "io"
	// "io/ioutil"
	// "os"
	// "path/filepath"
	"fmt"
	"strconv"
	"strings"

	// "gopkg.in/src-d/go-git.v4/utils/fs"
	// "strings"
	"github.com/wadahiro/gitss/server/config"
	"github.com/wadahiro/gitss/server/util"

	"bytes"
	// "io"
	"log"
	"net/http"

	gitm "github.com/gogits/git-module"
	// "gopkg.in/src-d/go-git.v4"
	// core "gopkg.in/src-d/go-git.v4/core"
	"github.com/pkg/errors"
)

type GitRepoReader struct {
	config *config.Config
}

func NewGitRepoReader(config *config.Config) *GitRepoReader {
	reader := &GitRepoReader{config: config}
	return reader
}

func GetRepoNameFromUrl(url string) string {
	splitedUrl := strings.Split(url, "/")
	repoName := splitedUrl[len(splitedUrl)-1]
	// Drop ".git" from repoName
	splitedRepoNames := strings.Split(repoName, ".git")
	repoName = splitedRepoNames[0]
	return repoName
}

func (r *GitRepoReader) CloneGitRepo(organization string, project string, url string) (*GitRepo, error) {
	repoName := GetRepoNameFromUrl(url)
	gitRepoPath := getGitRepoPath(r.config.GitDataDir, organization, project, repoName)

	err := gitm.Clone(url, gitRepoPath,
		gitm.CloneRepoOptions{Mirror: true})

	if err != nil {
		fmt.Println(err)
	}

	return r.GetGitRepo(organization, project, repoName)
}

func (r *GitRepoReader) GetGitRepo(organization string, project string, repoName string) (*GitRepo, error) {
	gitRepoPath := getGitRepoPath(r.config.GitDataDir, organization, project, repoName)

	repo, err := NewGitRepo(organization, project, repoName, gitRepoPath, r.config)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

type GitRepo struct {
	Organization string
	Project      string
	Repository   string
	Path         string
	gitmRepo     *gitm.Repository
	Config       *config.Config
}

type Source struct {
	Offset  int    `json:"offset"`
	Preview string `json:"preview"`
	Hits    []int  `json:"hits"`
}

func NewGitRepo(organization string, projectName string, repoName string, repoPath string, config *config.Config) (*GitRepo, error) {
	gitmRepo, err := gitm.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}

	return &GitRepo{Organization: organization, Project: projectName, Repository: repoName, Path: repoPath, gitmRepo: gitmRepo, Config: config}, nil
}

func (r *GitRepo) FetchAll() error {
	cmd := gitm.NewCommand("fetch")
	cmd.AddArguments("--all")
	cmd.AddArguments("--prune")

	_, err := cmd.RunInDirTimeout(-1, r.Path)
	return err
}

func (r *GitRepo) GetBranches() ([]string, error) {
	b, err := r.gitmRepo.GetBranches()
	if err != nil {
		return nil, errors.Wrapf(err, `Failed to get branch list. cmd: "show-ref --heads"`)
	}
	return b, nil
}

func (r *GitRepo) GetLatestCommitIdsMap() (map[string]string, map[string]string, error) {
	branches, err := r.GetBranches()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Failed to get branch list")
	}

	tags, err := r.gitmRepo.GetTags()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Failed to get tag list")
	}

	// master -> refs/heads/master
	b := []string{}
	for _, branch := range branches {
		b = append(b, gitm.BRANCH_PREFIX+branch)
	}

	// v1.0 => refs/heads/
	for _, tag := range tags {
		b = append(b, gitm.TAG_PREFIX+tag)
	}

	// get commitIds
	stdout, err := gitm.NewCommand("show-ref", "--verify").AddArguments(b...).RunInDir(r.Path)
	if err != nil {
		return nil, nil, errors.Wrapf(err, `Failed to get branch commitIds. cmd: "show-ref --verify %s"`, strings.Join(b, " "))
	}

	resultBranches := make(map[string]string)
	resultTags := make(map[string]string)

	rows := strings.Split(stdout, "\n")
	for _, row := range rows[:len(rows)-1] {
		columns := strings.Split(row, " ")
		if len(columns) != 2 {
			return nil, nil, errors.Errorf(`Already removed Ref? "git show-ref --verify %s" response: %s`, strings.Join(b, " "), row)
		}
		if strings.Contains(columns[1], gitm.BRANCH_PREFIX) {
			// refs/heads/master -> master
			branch := columns[1][len(gitm.BRANCH_PREFIX):]
			resultBranches[branch] = columns[0] // columns[0] is commitId
		} else {
			// refs/tags/v1.0 -> v1.0
			tag := columns[1][len(gitm.TAG_PREFIX):]
			resultTags[tag] = columns[0] // columns[0] is commitId
		}
	}

	return resultBranches, resultTags, nil
}

func (r *GitRepo) GetBranchCommitID(name string) (string, error) {
	return r.gitmRepo.GetBranchCommitID(name)
}

func (r *GitRepo) GetContainsBranches(commitId string) ([]string, error) {
	cmd := gitm.NewCommand("branch")
	cmd.AddArguments("--contains", commitId)

	stdout, err := cmd.RunInDir(r.Path)
	if err != nil {
		return nil, err
	}

	infos := strings.Split(stdout, "\n")

	branches := make([]string, len(infos)-1)
	for i, info := range infos[:len(infos)-1] {
		branches[i] = info[2:]
	}
	return branches, nil
}

func (r *GitRepo) GetBlobSize(blob string) (int64, error) {
	s, err := gitm.NewCommand("cat-file", "-s", blob).RunInDir(r.Path)
	if err != nil {
		return -1, err
	}
	s = strings.TrimRight(s, "\n")
	return strconv.ParseInt(s, 10, 64)
}

func (r *GitRepo) GetBlobContent(blob string) ([]byte, error) {
	b, err := gitm.NewCommand("cat-file", "-p", blob).RunInDirBytes(r.Path)
	if err != nil {
		return nil, err
	}
	b = bytes.TrimRight(b, "\n")
	return b, nil
}

// type ContentTypeWriter struct {
// 	pos   int
// 	bytes [512]byte
// }

// func (w *ContentTypeWriter) Write(p []byte) (n int, err error) {
// 	end := w.pos + len(p)
// 	if end > 512 {
// 		end = 512
// 	}

// 	copy(w.bytes[w.pos:end], p)
// 	w.pos = end

// 	return len(p), nil
// }

// func (r *GitRepo) DetectBlobContentType2(blob string) (string, error) {
// 	// Only the first 512 bytes are used to sniff the content type.
// 	stdout := new(ContentTypeWriter)
// 	stderr := new(bytes.Buffer)
// 	err := gitm.NewCommand("cat-file", "-p", blob).RunInDirPipeline(r.Path, stdout, stderr)
// 	if err != nil {
// 		return "", err
// 	}

// 	if r.config.Debug {
// 		// fmt.Println("ContentType size:", len(string(stdout.bytes[:])))
// 	}

// 	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
// 	contentType := http.DetectContentType(stdout.bytes[:])

// 	return contentType, nil
// }

func (r *GitRepo) DetectBlobContentType(blob string) (string, []byte, error) {
	b, err := gitm.NewCommand("cat-file", "-p", blob).RunInDirBytes(r.Path)
	if err != nil {
		return "", nil, err
	}

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(b)

	return contentType, b, nil
}

func (r *GitRepo) FilterBlob(blobId string, filter func(line string) bool, before int, after int) []util.TextPreview {
	b, _ := r.GetBlobContent(blobId)
	reader := strings.NewReader(string(b))

	previews := util.FilterTextPreview(reader, filter, before, after)

	return previews
}

type FileEntry struct {
	Blob string
	Path string
	Size int64
}

func (r *GitRepo) GetFileEntriesIterator(commitId string, callback func(fileEntry FileEntry)) error {
	// see https://git-scm.com/docs/git-ls-tree
	s, err := gitm.NewCommand("ls-tree", "-r", "-l", "--abbrev=40", commitId).RunInDir(r.Path)
	if err != nil {
		return errors.Wrapf(err, `Faild to get file list. cmd: "git ls-tree -r -l --abbrev=40 %s"`, commitId)
	}
	s = strings.TrimRight(s, "\n")
	rows := strings.Split(s, "\n")

	for i := range rows {
		row := rows[i]
		pathColumns := strings.Split(row, "\t")
		columns := strings.Fields(row)

		if len(pathColumns) != 2 {
			return errors.Errorf("Unexpected git ls-tree output. %s" + row)
		}

		blob := columns[2]
		size, _ := strconv.ParseInt(columns[3], 10, 64)

		path := strings.Split(row, "\t")[1]

		f := FileEntry{Blob: blob, Size: size, Path: path}

		callback(f)
	}

	return nil
}

func (r *GitRepo) GetFileEntries(commitId string) ([]FileEntry, error) {
	list := []FileEntry{}

	err := r.GetFileEntriesIterator(commitId, func(f FileEntry) {
		list = append(list, f)
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

type GitFile struct {
	Locations map[string]GitFileLocation
	Size      int64
}

type GitFileLocation struct {
	Branches []string
	Tags     []string
}

func (r *GitRepo) GetFileEntriesMap(excludeBranches []string) (map[string]GitFile, error) {
	branchesMap, tagsMap, err := r.GetLatestCommitIdsMap()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get branches/tags commitIds.")
	}

	files := make(map[string]GitFile)

	for branch, commitId := range branchesMap {
		entries, err := r.GetFileEntries(commitId)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get file entries of branch commit %s.", commitId)
		}

		for _, entry := range entries {
			_, ok := files[entry.Blob]
			if !ok {
				locations := make(map[string]GitFileLocation)
				locations[entry.Path] = GitFileLocation{Branches: []string{branch}, Tags: []string{}}
				files[entry.Blob] = GitFile{Locations: locations, Size: entry.Size}
			} else {
				_, ok := files[entry.Blob].Locations[entry.Path]
				if !ok {
					files[entry.Blob].Locations[entry.Path] = GitFileLocation{Branches: []string{branch}, Tags: []string{}}
				} else {
					location := files[entry.Blob].Locations[entry.Path]
					location.Branches = append(location.Branches, branch)
					files[entry.Blob].Locations[entry.Path] = location
				}
			}
		}
	}

	for tag, commitId := range tagsMap {
		entries, err := r.GetFileEntries(commitId)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get file entries of tag commit %s.", commitId)
		}

		for _, entry := range entries {
			_, ok := files[entry.Blob]
			if !ok {
				locations := make(map[string]GitFileLocation)
				locations[entry.Path] = GitFileLocation{Branches: []string{}, Tags: []string{tag}}
				files[entry.Blob] = GitFile{Locations: locations, Size: entry.Size}
			} else {
				_, ok := files[entry.Blob].Locations[entry.Path]
				if !ok {
					files[entry.Blob].Locations[entry.Path] = GitFileLocation{Branches: []string{}, Tags: []string{tag}}
				} else {
					location := files[entry.Blob].Locations[entry.Path]
					location.Tags = append(location.Tags, tag)
					files[entry.Blob].Locations[entry.Path] = location
				}
			}
		}
	}
	return files, nil
}

func (r *GitRepo) GetDiffEntriesIterator(from string, to string, callback func(fileEntry FileEntry, status string)) error {
	// see https://git-scm.com/docs/diff
	stdout, err := gitm.NewCommand("diff", "--raw", "--abbrev=40", "-z", from, to).RunInDirTimeout(-1, r.Path)
	if err != nil {
		return err
	}

	parts := bytes.Split(stdout, []byte{'\u0000'})

	// fmt.Println(len(parts))

	for index := 0; index < len(parts); index++ {
		row := parts[index]
		cols := bytes.Fields(row)

		if len(cols) <= 1 {
			continue
		}

		// fmt.Println("col", len(cols))

		oldBlob := string(cols[2])
		newBlob := string(cols[3])
		status := string(cols[4])[0:1]

		// fmt.Println(string(row))
		// fmt.Println(oldBlob, newBlob, status)

		// See 'Possible status letters' in https://git-scm.com/docs/git-diff
		switch status {
		case "A": // Add case
			index++
			path := string(parts[index])
			callback(FileEntry{Blob: newBlob, Path: path}, "A")

		case "C": // Copy case
			index = index + 2
			toPath := string(parts[index])
			callback(FileEntry{Blob: newBlob, Path: toPath}, "A")

		case "D": // Delete case
			index++
			path := string(parts[index])
			callback(FileEntry{Blob: oldBlob, Path: path}, "D")

		case "M", "T": // Modify case and change in the type of the file case
			index++
			if oldBlob != newBlob {
				path := string(parts[index])
				callback(FileEntry{Blob: oldBlob, Path: path}, "D")
				callback(FileEntry{Blob: newBlob, Path: path}, "A")
			}

		case "R": // Rename case
			index++
			fromPath := string(parts[index])
			index++
			toPath := string(parts[index])
			callback(FileEntry{Blob: oldBlob, Path: fromPath}, "D")
			callback(FileEntry{Blob: newBlob, Path: toPath}, "A")

		case "U": // Unmerge case
			continue
		case "X": // Delete case
			log.Println("unknown change type (most probably a git bug, please report it)")
		}
	}
	return nil
}

func (r *GitRepo) GetDiffList(from string, to string) ([]FileEntry, []FileEntry, error) {
	stdout, err := gitm.NewCommand("diff", "--raw", "--abbrev=40", "-z", from, to).RunInDirTimeout(-1, r.Path)
	if err != nil {
		return nil, nil, err
	}

	addList := []FileEntry{}
	delList := []FileEntry{}

	parts := bytes.Split(stdout, []byte{'\u0000'})

	// fmt.Println(len(parts))

	for index := 0; index < len(parts); index++ {
		row := parts[index]
		cols := bytes.Fields(row)

		if len(cols) <= 1 {
			continue
		}

		// fmt.Println("col", len(cols))

		oldBlob := string(cols[2])
		newBlob := string(cols[3])
		status := string(cols[4])[0:1]

		// fmt.Println(string(row))
		// fmt.Println(oldBlob, newBlob, status)

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

		case "M", "T": // Modify case and change in the type of the file case
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

func (r *GitRepo) ExistsInCommit(commitId string, filePath string, blobId string) (bool, error) {
	// see https://git-scm.com/docs/git-ls-tree
	s, err := gitm.NewCommand("ls-tree", "--abbrev=40", commitId, "--", filePath).RunInDir(r.Path)
	if err != nil {
		return false, err
	}

	// log.Println(s)

	result := strings.Split(s, "\n")
	if len(result) < 2 {
		// Not found case
		return false, nil
	}

	columns := strings.Fields(result[0])
	if len(columns) != 4 {
		return false, errors.Errorf("Unexpected ls-tree response. command: git ls-tree -l --abbrev=40 %s -- %s response: %s", commitId, filePath, s)
	}

	return columns[2] == blobId, nil
}

func getGitRepoPath(GitDataDir string, organization string, project string, repoName string) string {
	repoPath := fmt.Sprintf("%s/%s/%s/%s.git", GitDataDir, organization, project, repoName)
	return repoPath
}
