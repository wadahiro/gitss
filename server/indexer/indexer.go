package indexer

import (
	"github.com/wadahiro/GitS/server/util"
)

type Indexer interface {
	CreateFileIndex(project string, repo string, branch string, fileName string, blob string, content string) error
	UpsertFileIndex(project string, repo string, branch string, fileName string, blob string, content string) error
	SearchQuery(query string) SearchResult
}

type Metadata struct {
	Project string `json:"project"`
	Repo    string `json:"repo"`
	Refs    string `json:"refs"`
	Path    string `json:"path"`
	Ext     string `json:"ext"`
}

type SearchResult struct {
	Time       float64 `json:"time"`
	Size       int64     `json:"size"`
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
