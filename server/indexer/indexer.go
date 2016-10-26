package indexer

import (
	"regexp"
	// "log"
	"fmt"
	"path"

	"github.com/wadahiro/gitss/server/repo"
	"github.com/wadahiro/gitss/server/util"
)

type Indexer interface {
	CreateFileIndex(requestFileIndex FileIndex) error
	UpsertFileIndex(requestFileIndex FileIndex) error
	BatchFileIndex(operations []FileIndexOperation) error

	SearchQuery(query string) SearchResult
}

type BatchMethod int

const (
	ADD BatchMethod = iota
	DELETE
)

type FileIndexOperation struct {
	Method    BatchMethod
	FileIndex FileIndex
}

type FileIndex struct {
	Blob     string   `json:"blob"`
	Metadata Metadata `json:"metadata"`
	Content  string   `json:"content"`
}

type Metadata struct {
	Organization string   `json:"organization"`
	Project      string   `json:"project"`
	Repository   string   `json:"repository"`
	Refs         []string `json:"refs"`
	Path         string   `json:"path"`
	Ext          string   `json:"ext"`
}

type SearchResult struct {
	Time       float64 `json:"time"`
	Size       int64   `json:"size"`
	Limit      int     `json:"limit"`
	isLastPage bool    `json:"isLastPage"`
	Current    int     `json:"current"`
	Next       int     `json:"next"`
	Hits       []Hit   `json:"hits"`
}

type Hit struct {
	Source Source `json:"_source"`
	// Highlight map[string][]string `json:"highlight"`
	Preview []util.TextPreview `json:"preview"`
}

type Source struct {
	Blob     string   `json:"blob"`
	Metadata Metadata `json:"metadata"`
}

type HighlightSource struct {
	Offset  int    `json:"offset"`
	Content string `json:"content"`
}

func getGitRepo(reader *repo.GitRepoReader, s *Source) (*repo.GitRepo, error) {
	repo, err := reader.GetGitRepo(s.Metadata.Organization, s.Metadata.Project, s.Metadata.Repository)
	return repo, err
}

func NewFileIndex(blob string, organization string, project string, repo string, ref string, path string, content string) FileIndex {
	fileIndex := FileIndex{
		Blob:    blob,
		Content: content,
		Metadata: Metadata{
			Organization: organization,
			Project:      project,
			Repository:   repo,
			Refs:         []string{ref},
			Path:         path,
		},
	}
	return fileIndex
}

func fillFileExt(fileIndex *FileIndex) {
	ext := path.Ext(fileIndex.Metadata.Path)
	fileIndex.Metadata.Ext = ext
}

func find(refs []string, f func(ref string, i int) bool) string {
	for index, x := range refs {
		if f(x, index) == true {
			return x
		}
	}
	return ""
}

func filter(f func(s Metadata, i int) bool, s []Metadata) []Metadata {
	ans := make([]Metadata, 0)
	for index, x := range s {
		if f(x, index) == true {
			ans = append(ans, x)
		}
	}
	return ans
}

func getDocId(fileIndex FileIndex) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", fileIndex.Metadata.Organization, fileIndex.Metadata.Project, fileIndex.Metadata.Repository, fileIndex.Blob, fileIndex.Metadata.Path)
}

func mergeRef(fileIndex *FileIndex, refs []string) bool {
	addRef := []string{}

	for i := range refs {
		found := find(fileIndex.Metadata.Refs, func(x string, j int) bool {
			return refs[i] == x
		})
		if found == "" {
			addRef = append(addRef, refs[i])
		}
	}

	// Same case
	if len(addRef) == 0 {
		return true
	}

	// Add case
	fileIndex.Metadata.Refs = append(fileIndex.Metadata.Refs, addRef...)

	return false
}

func mergeSet(m1, m2 map[string]struct{}) map[string]struct{} {
	ans := make(map[string]struct{})

	for k, v := range m1 {
		ans[k] = v
	}
	for k, v := range m2 {
		ans[k] = v
	}
	return (ans)
}

func getHitWords(hitTag *regexp.Regexp, contents []string) map[string]struct{} {
	hitWordsSet := make(map[string]struct{})

	for _, content := range contents {
		groups := hitTag.FindAllStringSubmatch(content, -1)
		// log.Println("hit", len(groups))
		for _, group := range groups {
			for i, g := range group {
				if i == 0 {
					continue
				}
				// log.Println("hit2", g)
				hitWordsSet[g] = struct{}{}
			}
		}
	}

	return hitWordsSet
}
