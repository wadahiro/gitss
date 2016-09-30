package indexer

type Indexer interface {
	CreateFileIndex(project string, repo string, branch string, fileName string, blob string, content string)
}
