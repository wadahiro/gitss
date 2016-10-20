package controller

import (
	// "log"
	// "time"
	// "bytes"
	// "fmt"
	// "log"
	// "strings"

	"github.com/gin-gonic/gin"
	"github.com/wadahiro/gitss/server/indexer"
)

func SearchIndex(c *gin.Context) {
	indexer := getIndexer(c)

	query := c.Query("q")

	result := indexer.SearchQuery(query)

	c.JSON(200, result)
}

func getIndexer(c *gin.Context) indexer.Indexer {
	r, _ := c.Get("indexer")
	indexer := r.(indexer.Indexer)

	return indexer
}

func getGitDataDir(c *gin.Context) string {
	r, _ := c.Get("gitDataDir")
	gitDataDir := r.(string)

	return gitDataDir
}
