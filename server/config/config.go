package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	// "path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/codegangsta/cli"
)

var fileMutex sync.Mutex

type Config struct {
	DataDir     string
	GitDataDir  string
	ConfDir     string
	Port        int
	IndexerType string
	SizeLimit   int64
	Schedule    string
	Debug       bool
	settings    []OrganizationSetting
}

type LatestIndex struct {
	Organization string `json:"organization"`
	Project      string `json:"project"`
	Repository   string `json:"repository"`
	Ref          string `json:"ref"`
}

func NewConfig(c *cli.Context, debug bool) *Config {
	port := c.GlobalInt("port")
	dataDir := c.GlobalString("data")
	gitDataDir := dataDir + "/" + "git"
	confDir := dataDir + "/" + "conf"

	indexerType := c.GlobalString("indexer")

	sizeLimit := c.Int64("sizeLimit")

	schedule := c.GlobalString("schedule")

	config := &Config{
		DataDir:     dataDir,
		GitDataDir:  gitDataDir,
		ConfDir:     confDir,
		Port:        port,
		IndexerType: indexerType,
		SizeLimit:   sizeLimit,
		Schedule:    schedule,
		Debug:       debug,
	}

	config.init()

	return config
}

func (c *Config) init() {
	if err := os.MkdirAll(c.GitDataDir, 0644); err != nil {
		log.Fatalln(err)
	}
	if err := os.MkdirAll(c.ConfDir, 0644); err != nil {
		log.Fatalln(err)
	}
	c.refreshSettings()
}

func (c *Config) Sync() {
	c.refreshSettings()
}

func (c *Config) refreshSettings() error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	list := []OrganizationSetting{}

	files, err := filepath.Glob(c.ConfDir + "/*.json")
	if err != nil {
		return err
	}

	for _, path := range files {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("Not found config, probably deleted. %s\n", path) // NotFound
			continue
		}
		var organizationSetting OrganizationSetting
		json.Unmarshal(b, &organizationSetting)

		list = append(list, organizationSetting)
	}

	// cache update
	c.settings = list

	return nil
}

func (c *Config) GetSettings() []OrganizationSetting {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	return c.settings
}

type OrganizationSetting struct {
	Name     string           `json:"name"`
	Projects []ProjectSetting `json:"projects"`
}
type ProjectSetting struct {
	Name         string              `json:"name"`
	Repositories []RepositorySetting `json:"repositories"`
}
type RepositorySetting struct {
	Url  string       `json:"url"`
	name string       `json:"-"`
	Refs []RefSetting `json:"refs"`
}

func (r *RepositorySetting) GetName() string {
	if r.name != "" {
		return r.name
	}

	// cache
	r.name = repoUrlToName(r.Url)

	return r.name
}

func repoUrlToName(repositoryUrl string) string {
	url := strings.Split(repositoryUrl, "/")
	var repoName string
	if len(url) > 0 {
		repoName = url[len(url)-1]
		if strings.HasSuffix(strings.ToLower(repoName), ".git") {
			i := strings.LastIndex(repoName, ".")
			repoName = repoName[:i]
		}
	}
	return repoName
}

type RefSetting struct {
	Name   string `json:"name"`
	Latest string `json:"latest"`
}

func (c *Config) findOrganizationSetting(name string) (*OrganizationSetting, bool) {
	for i := range c.settings {
		if c.settings[i].Name == name {
			return &c.settings[i], true
		}
	}
	return &OrganizationSetting{}, false
}

func (c *Config) findProjectSetting(organization string, project string) (*ProjectSetting, bool) {
	organizationSetting, ok := c.findOrganizationSetting(organization)
	if ok {
		for i := range organizationSetting.Projects {
			if organizationSetting.Projects[i].Name == project {
				return &organizationSetting.Projects[i], true
			}
		}
	}
	return &ProjectSetting{}, false
}

func (c *Config) findRepositorySetting(organization string, project string, repository string) (*RepositorySetting, bool) {
	projectSetting, ok := c.findProjectSetting(organization, project)
	if ok {
		for i := range projectSetting.Repositories {
			if projectSetting.Repositories[i].GetName() == repository {
				return &projectSetting.Repositories[i], true
			}
		}
	}
	return &RepositorySetting{}, false
}

func (c *Config) findRefSetting(organization string, project string, repository string, ref string) (*RefSetting, bool) {
	repositorySetting, ok := c.findRepositorySetting(organization, project, repository)
	if ok {
		for i := range repositorySetting.Refs {
			refSetting := &repositorySetting.Refs[i]
			if refSetting.Name == ref {
				return refSetting, true
			}
		}
	}
	return &RefSetting{}, false
}

func (c *Config) GetRefs(organization string, project string, repository string) ([]RefSetting, bool) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	repositorySetting, ok := c.findRepositorySetting(organization, project, repository)
	if ok {
		return repositorySetting.Refs, true
	}

	return []RefSetting{}, false
}

func (c *Config) GetIndexedCommitID(organization string, project string, repository string, ref string) (string, bool) {
	refSetting, ok := c.findRefSetting(organization, project, repository, ref)
	if ok {
		return refSetting.Latest, true
	}
	return "", false
}

func (c *Config) AddRepositorySetting(organization string, project string, url string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	organizationSetting, ok := c.findOrganizationSetting(organization)
	if !ok {
		organizationSetting = &OrganizationSetting{
			Name: organization,
			Projects: []ProjectSetting{
				ProjectSetting{
					Name: project,
					Repositories: []RepositorySetting{
						RepositorySetting{
							Url:  url,
							Refs: []RefSetting{},
						},
					},
				},
			},
		}
		c.settings = append(c.settings, *organizationSetting)

	} else {
		projectSetting, ok := c.findProjectSetting(organization, project)
		if !ok {
			projectSetting = &ProjectSetting{
				Name: project,
				Repositories: []RepositorySetting{
					RepositorySetting{
						Url:  url,
						Refs: []RefSetting{},
					},
				},
			}
			organizationSetting.Projects = append(organizationSetting.Projects, *projectSetting)

		} else {
			repositorySetting, ok := c.findRepositorySetting(organization, project, repoUrlToName(url))
			if !ok {
				repositorySetting = &RepositorySetting{
					Url: url,
					Refs: []RefSetting{},
				}
				projectSetting.Repositories = append(projectSetting.Repositories, *repositorySetting)

			} else {
				return errors.Errorf("Already exists %s:%s/%s", organization, project, repoUrlToName(url))
			}
		}
	}

	// write!
	c.writeSetting(organization)

	return nil
}

func (c *Config) UpdateLatestIndex(url string, organization string, project string, repository string, ref string, commitId string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// update case
	refSetting, ok := c.findRefSetting(organization, project, repository, ref)
	if ok {
		refSetting.Latest = commitId
	} else {
		// add ref case
		repositorySetting, ok := c.findRepositorySetting(organization, project, repository)
		if ok {
			repositorySetting.Refs = append(repositorySetting.Refs, RefSetting{Name: ref, Latest: commitId})
		}
	}

	// write
	return c.writeSetting(organization)
}

func (c *Config) DeleteLatestIndexRefs(organization string, project string, repository string, removeRefs []string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	repositorySetting, ok := c.findRepositorySetting(organization, project, repository)
	if ok {
		newRefSettings := []RefSetting{}
		for i := range repositorySetting.Refs {
			ref := repositorySetting.Refs[i]

			found := false
			for _, removeRef := range removeRefs {
				if ref.Name == removeRef {
					found = true
					break
				}
			}
			if !found {
				newRefSettings = append(newRefSettings, repositorySetting.Refs[i])
			}
		}
		repositorySetting.Refs = newRefSettings

		return c.writeSetting(organization)
	}
	return errors.Errorf("Not found repository setting for %s:%s/%s", organization, project, repository)
}

func (c *Config) writeSetting(organization string) error {
	organizationSetting, ok := c.findOrganizationSetting(organization)
	if ok {
		content, _ := json.MarshalIndent(organizationSetting, "", "  ")
		fileName := fmt.Sprintf("%s/%s.json", c.ConfDir, organization)
		return ioutil.WriteFile(fileName, content, os.ModePerm)
	}
	return errors.Errorf("Not found organization setting for %s", organization)
}
