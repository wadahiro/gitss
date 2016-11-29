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

type IndexdListResult struct {
	Result []config.Indexed  `json:"result"`
}

func GetIndexedList(c *gin.Context) {
	cfg := getConfig(c)

	list := []config.Indexed{}

	settings := cfg.GetSettings()
	for _, setting:= range settings {
		projects := setting.GetProjects()
		for _ , projectSetting := range projects {
			for _ , repoSetting := range projectSetting.Repositories {
				indexed := cfg.GetIndexed(setting.GetName(), projectSetting.Name, repoSetting.GetName())
				
				list = append(list, indexed)
			}
		}
	}
	
	c.JSON(200, IndexdListResult{
		Result: list,
	})
}