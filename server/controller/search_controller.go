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
		exts, _ := c.Request.Form["x"]
		organizations, _ := c.Request.Form["o"]
		projects, _ := c.Request.Form["p"]
		repositories, _ := c.Request.Form["r"]
		branches, _ := c.Request.Form["b"]
		result := i.SearchQuery(q[0], indexer.FilterParams{Exts: exts, Organizations: organizations, Projects: projects, Repositories: repositories, Refs: branches})

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
