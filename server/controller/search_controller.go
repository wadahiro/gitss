package controller

import (
	// "log"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/wadahiro/gits/server/indexer"
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
