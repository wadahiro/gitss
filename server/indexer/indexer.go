package indexer
import (
	// "encoding/json"
)
type Indexer interface {
	CreateFileIndex(project string, repo string, branch string, fileName string, blob string, content string) error
	UpsertFileIndex(project string, repo string, branch string, fileName string, blob string, content string) error
	SearchQuery(query string) []Hit
}
