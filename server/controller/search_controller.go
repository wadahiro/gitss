package controller

import (
	// "log"
	// "time"
	// "bytes"
	// "fmt"
	// "log"

	"github.com/gin-gonic/gin"
	"github.com/wadahiro/gitss/server/indexer"
)

func SearchIndex(c *gin.Context) {
	i := getIndexer(c)

	c.Request.ParseForm()

	q, ok := c.Request.Form["q"]
	// fmt.Println(q, ok)
	if ok {
		ext, ok := c.Request.Form["ext"]
		if !ok {
			ext = []string{}
		}
		result := i.SearchQuery(q[0], indexer.FilterParams{Ext: ext})

		if result.FilterParams.Ext == nil {
			result.FilterParams.Ext = []string{}
		}

		c.JSON(200, result)
	} else {
		c.JSON(200, indexer.SearchResult{})
	}
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
