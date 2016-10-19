package indexer

import (
	"bytes"
	"regexp"
	// "log"

	"github.com/wadahiro/gitss/server/repo"
	"github.com/wadahiro/gitss/server/util"
)

type Indexer interface {
	CreateFileIndex(organization string, project string, repo string, branch string, fileName string, blob string, content string) error
	UpsertFileIndex(organization string, project string, repo string, branch string, fileName string, blob string, content string) error
	SearchQuery(query string) SearchResult
}

type FileIndex struct {
	Blob     string     `json:"blob"`
	Metadata []Metadata `json:"metadata"`
	Content  string     `json:"content"`
}

type Metadata struct {
	Organization string `json:"organization"`
	Project      string `json:"project"`
	Repository   string `json:"repository"`
	Refs         string `json:"refs"`
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

func getFileContent(repo *repo.GitRepo, s *Source) string {
	blob, _ := repo.GetBlob(s.Blob)

	r, _ := blob.Reader()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	text := buf.String()

	return text
}

func find(f func(s Metadata, i int) bool, s []Metadata) *Metadata {
	for index, x := range s {
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

func mergeFileIndex(fileIndex *FileIndex, organization string, project string, repo string, refs string, filePath string, ext string) {
	f := func(x Metadata, i int) bool {
		return x.Organization == organization &&
			x.Project == project &&
			x.Repository == repo &&
			x.Refs == refs &&
			x.Path == filePath
	}
	found := find(f, fileIndex.Metadata)
	// log.Println("before:", fileIndex.Metadata)
	if found == nil {
		fileIndex.Metadata = append(fileIndex.Metadata, Metadata{Organization: organization, Project: project, Repository: repo, Refs: refs, Path: filePath, Ext: ext})
	}
	// log.Println("merged:", fileIndex.Metadata)
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
