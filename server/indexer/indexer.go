package indexer

import (
	"regexp"
	// "log"
	"path"

	"github.com/wadahiro/gitss/server/repo"
	"github.com/wadahiro/gitss/server/util"
)

type Indexer interface {
	CreateFileIndex(requestFileIndex FileIndex) error
	UpsertFileIndex(requestFileIndex FileIndex) error
	BatchFileIndex(requestFileIndex []FileIndex, mode BatchMethod) error

	SearchQuery(query string) SearchResult
}

type BatchMethod int

const (
    ADD BatchMethod = iota
    DELETE
)

type FileIndex struct {
	Blob     string     `json:"blob"`
	Metadata []Metadata `json:"metadata"`
	Content  string     `json:"content"`
}

type Metadata struct {
	Organization string `json:"organization"`
	Project      string `json:"project"`
	Repository   string `json:"repository"`
	Ref          string `json:"ref"`
	Path         string `json:"path"`
	Ext          string `json:"ext"`
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
	Blob     string     `json:"blob"`
	Metadata []Metadata `json:"metadata"`
}

type HighlightSource struct {
	Offset  int    `json:"offset"`
	Content string `json:"content"`
}

func getGitRepo(reader *repo.GitRepoReader, s *Source) *repo.GitRepo {
	repo := reader.GetGitRepo(s.Metadata[0].Organization, s.Metadata[0].Project, s.Metadata[0].Repository)
	return repo
}

// func getFileContent(repo *repo.GitRepo, s *Source) string {
// 	blob, _ := repo.GetBlob(s.Blob)

// 	r, _ := blob.Reader()

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r)
// 	text := buf.String()

// 	return text
// }

func NewFileIndex(blob string, organization string, project string, repo string, ref string, path string, content string) FileIndex {
	fileIndex := FileIndex{
		Blob:    blob,
		Content: content,
		Metadata: []Metadata{
			Metadata{
				Organization: organization,
				Project:      project,
				Repository:   repo,
				Ref:          ref,
				Path:         path,
			},
		},
	}
	return fileIndex
}

func fillFileExt(fileIndex *FileIndex) {
	for i := range fileIndex.Metadata {
		ext := path.Ext(fileIndex.Metadata[i].Path)
		fileIndex.Metadata[i].Ext = ext
	}
}

func find(mList []Metadata, f func(m Metadata, i int) bool) *Metadata {
	for index, x := range mList {
		if f(x, index) == true {
			return &x
		}
	}
	return nil
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

func mergeFileIndex(fileIndex *FileIndex, metadata []Metadata) bool {
	addMetadata := []Metadata{}

	for i := range metadata {
		m := metadata[i]

		found := find(fileIndex.Metadata, func(n Metadata, j int) bool {
			return m.Organization == n.Organization &&
				m.Project == n.Project &&
				m.Repository == n.Repository &&
				m.Ref == n.Ref &&
				m.Path == n.Path
		})
		if found == nil {
			addMetadata = append(addMetadata, m)
		}
	}

	// Same metadata case
	if len(addMetadata) == 0 {
		return true
	}

	// Add metadata case
	fileIndex.Metadata = append(fileIndex.Metadata, addMetadata...)

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
