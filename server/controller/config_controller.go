package controller

import (
	// "log"
	// "time"
	// "bytes"
	// "fmt"
	// "log"

	"github.com/gin-gonic/gin"
	"github.com/wadahiro/gitss/server/config"
)

type IndexdStatisticsResult struct {
	Count   Count            `json:"count"`
	Indexes []config.Indexed `json:"indexes"`
}

type Count struct {
	Organization int    `json:"organization"`
	Project      int    `json:"project"`
	Repository   int    `json:"repository"`
	Branch       int    `json:"branch"`
	Tag          int    `json:"tag"`
	Document     uint64 `json:"document"`
}

func GetIndexStatistics(c *gin.Context) {
	cfg := getConfig(c)
	i := getIndexer(c)

	list := []config.Indexed{}
	projectCount := 0
	repositoryCount := 0
	branchCount := 0
	tagCount := 0

	settings := cfg.GetSettings()
	for _, setting := range settings {
		projects := setting.GetProjects()
		projectCount += len(projects)
		for _, projectSetting := range projects {
			repositoryCount += len(projectSetting.Repositories)
			for _, repoSetting := range projectSetting.Repositories {
				indexed := cfg.GetIndexed(setting.GetName(), projectSetting.Name, repoSetting.GetName())

				list = append(list, indexed)

				branchCount += len(indexed.Branches)
				tagCount += len(indexed.Tags)
			}
		}
	}

	docCount, err := i.Count()
	if err != nil {
		errorJson := make(map[string]string)
		errorJson["error"] = "Cannot get statistics"
		c.JSON(500, errorJson)
	}

	c.JSON(200, IndexdStatisticsResult{
		Count: Count{
			Organization: len(settings),
			Project:      projectCount,
			Repository:   repositoryCount,
			Branch:       branchCount,
			Tag:          tagCount,
			Document:     docCount,
		},
		Indexes: list,
	})
}
