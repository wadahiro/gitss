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
	DeleteIndexByRefs(organization string, project string, repository string, refs []string) error

	SearchQuery(query string, filters FilterParams, page int) SearchResult
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
	Metadata
	FullRefs []string `json:"fullRefs"`
	Content  string   `json:"content"`
}

type Metadata struct {
	Blob         string   `json:"blob"`
	Organization string   `json:"organization"`
	Project      string   `json:"project"`
	Repository   string   `json:"repository"`
	Refs         []string `json:"refs"`
	Path         string   `json:"path"`
	Ext          string   `json:"ext"`
	Size         int64    `json:"size"`
}

type SearchResult struct {
	Query         string              `json:"query"`
	FilterParams  FilterParams        `json:"filterParams"`
	Time          float64             `json:"time"`
	Size          int64               `json:"size"`
	Limit         int                 `json:"limit"`
	isLastPage    bool                `json:"isLastPage"`
	Current       int                 `json:"current"`
	Next          int                 `json:"next"`
	Hits          []Hit               `json:"hits"`
	FullRefsFacet []OrganizationFacet `json:"fullRefsFacet"`
	Facets        FacetResults        `json:"facets"`
}

type OrganizationFacet struct {
	Term     string         `json:"term"`
	Count    int            `json:"count"`
	Projects []ProjectFacet `json:"projects"`
}
type ProjectFacet struct {
	Term         string            `json:"term"`
	Count        int               `json:"count"`
	Repositories []RepositoryFacet `json:"repositories"`
}
type RepositoryFacet struct {
	Term  string     `json:"term"`
	Count int        `json:"count"`
	Refs  []RefFacet `json:"refs"`
}
type RefFacet struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}

type Hit struct {
	Metadata
	Keyword []string `json:"keyword"`
	// Highlight map[string][]string `json:"highlight"`
	Preview []util.TextPreview `json:"preview"`
}

type HighlightSource struct {
	Offset  int    `json:"offset"`
	Content string `json:"content"`
}

type FacetResults map[string]FacetResult

type FacetResult struct {
	Field   string     `json:"field"`
	Total   int        `json:"total"`
	Missing int        `json:"missing"`
	Other   int        `json:"other"`
	Terms   TermFacets `json:"terms"`
}

type TermFacets []TermFacet

type TermFacet struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}

type FilterParams struct {
	Exts          []string `json:"x,omitempty"`
	Organizations []string `json:"o,omitempty"`
	Projects      []string `json:"p,omitempty"`
	Repositories  []string `json:"r,omitempty"`
	Refs          []string `json:"b,omitempty"`
}

func getGitRepo(reader *repo.GitRepoReader, fileIndex *FileIndex) (*repo.GitRepo, error) {
	repo, err := reader.GetGitRepo(fileIndex.Organization, fileIndex.Project, fileIndex.Repository)
	return repo, err
}

func NewFileIndex(blob string, organization string, project string, repo string, ref string, path string, content string) FileIndex {
	fileIndex := FileIndex{
		Metadata: Metadata{
			Blob:         blob,
			Organization: organization,
			Project:      project,
			Repository:   repo,
			Refs:         []string{ref},
			Path:         path,
		},
		Content: content,
	}
	return fileIndex
}

const NO_EXT = "/noext/"

func GetExt(p string) string {
	ext := path.Ext(p)
	if ext == "" {
		return NO_EXT
	}
	return ext
}

func fillFileIndex(fileIndex *FileIndex) {
	// ext
	ext := GetExt(fileIndex.Metadata.Path)
	fileIndex.Metadata.Ext = ext

	// full_refs
	fullRefs := make([]string, 0, len(fileIndex.Metadata.Refs))
	for _, ref := range fileIndex.Metadata.Refs {
		fullRefs = append(fullRefs, fileIndex.Metadata.Organization+":"+fileIndex.Metadata.Project+"/"+fileIndex.Metadata.Repository+":"+ref)
	}
	fileIndex.FullRefs = fullRefs
}

func find(refs []string, f func(ref string, i int) bool) (string, bool) {
	for index, x := range refs {
		if f(x, index) == true {
			return x, true
		}
	}
	return "", false
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

func getDocId(fileIndex *FileIndex) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", fileIndex.Metadata.Organization, fileIndex.Metadata.Project, fileIndex.Metadata.Repository, fileIndex.Blob, fileIndex.Metadata.Path)
}

func mergeRef(fileIndex *FileIndex, refs []string) bool {
	addRefs := []string{}
	currentRefs := fileIndex.Metadata.Refs

	for i := range refs {
		_, found := find(currentRefs, func(x string, j int) bool {
			return refs[i] == x
		})
		if !found {
			addRefs = append(addRefs, refs[i])
		}
	}

	// Same case
	if len(addRefs) == 0 {
		return true
	}

	addFullRefs := make([]string, 0, len(addRefs))
	for _, ref := range addRefs {
		fullRef := fileIndex.Metadata.Organization + ":" + fileIndex.Metadata.Project + "/" + fileIndex.Metadata.Repository + ":" + ref
		addFullRefs = append(addFullRefs, fullRef)
	}

	// Add case
	fileIndex.Metadata.Refs = append(fileIndex.Metadata.Refs, addRefs...)
	fileIndex.FullRefs = append(fileIndex.FullRefs, addFullRefs...)

	return false
}

func removeRef(fileIndex *FileIndex, refs []string) bool {
	newRefs := []string{}
	currentRefs := fileIndex.Metadata.Refs

	for i := range currentRefs {
		_, found := find(refs, func(x string, j int) bool {
			return currentRefs[i] == x
		})
		if !found {
			newRefs = append(newRefs, currentRefs[i])
		}
	}

	// All delete case
	if len(newRefs) == 0 {
		return true
	}

	// Update case
	fileIndex.Metadata.Refs = newRefs

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
