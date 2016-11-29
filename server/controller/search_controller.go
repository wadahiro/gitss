package controller

import (
	// "log"
	// "time"
	// "bytes"
	// "fmt"
	// "log"
	"strconv"

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
		tags, _ := c.Request.Form["t"]

		reqPage, ok := c.Request.Form["i"]
		page := 0
		if ok {
			p, err := strconv.Atoi(reqPage[0])
			if err == nil {
				page = p
			}
		}

		result, err := i.SearchQuery(q[0], indexer.FilterParams{Exts: exts, Organizations: organizations, Projects: projects, Repositories: repositories, Branches: branches, Tags: tags}, page)

		if err != nil {
			c.AbortWithError(500, err)
			return
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
