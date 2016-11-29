package controller

import (
	// "log"
	// "time"
	// "bytes"
	// "fmt"
	// "log"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/wadahiro/gitss/server/config"
)

type FilterResult struct {
	Organizations []string `json:"organizations,omitempty"`
	Projects      []string `json:"projects,omitempty"`
	Repositories  []string `json:"repositories,omitempty"`
	Branches      []string `json:"branches,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

func GetBaseFilters(c *gin.Context) {
	cfg := getConfig(c)

	organization := c.Param("organization")
	project := c.Param("project")
	repository := c.Param("repository")

	if organization == "" {
		settings := cfg.GetSettings()

		organizations := []string{}
		for i := range settings {
			organizations = append(organizations, settings[i].GetName())
		}

		c.JSON(200, FilterResult{Organizations: organizations})
		return

	} else if project == "" {
		setting, ok := cfg.FindSetting(organization)
		if ok {
			projects := []string{}
			for _, projectSetting := range setting.GetProjects() {
				projects = append(projects, projectSetting.Name)
			}

			c.JSON(200, FilterResult{Organizations: []string{organization}, Projects: projects})
			return
		}

	} else if repository == "" {
		setting, ok := cfg.FindSetting(organization)
		if ok {
			projectSetting, ok := setting.FindProjectSetting(project)
			if ok {
				repositories := []string{}
				for _, repository := range projectSetting.Repositories {
					repositories = append(repositories, repository.GetName())
				}

				c.JSON(200, FilterResult{Organizations: []string{organization}, Projects: []string{project}, Repositories: repositories})
				return
			}
		}
	} else {
		indexed := cfg.GetIndexed(organization, project, repository)
		branches := []string{}
		for branch := range indexed.Branches {
			branches = append(branches, branch)
		}
		sort.Sort(sort.StringSlice(branches))

		tags := []string{}
		for tag := range indexed.Tags {
			tags = append(tags, tag)
		}
		sort.Sort(sort.StringSlice(tags))

		c.JSON(200, FilterResult{
			Organizations: []string{organization},
			Projects:      []string{project},
			Repositories:  []string{repository},
			Branches:      branches,
			Tags:          tags,
		})
		return
	}

	errorJson := make(map[string]string)
	errorJson["error"] = "Not found"
	c.JSON(404, errorJson)
}

func getConfig(c *gin.Context) *config.Config {
	r, _ := c.Get("config")
	config := r.(*config.Config)

	return config
}
