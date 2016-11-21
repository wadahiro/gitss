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
var indexedFileMutex sync.Mutex

type Config struct {
	DataDir     string
	GitDataDir  string
	ConfDir     string
	IndexedDir  string
	Port        int
	IndexerType string
	SizeLimit   int64
	Schedule    string
	Debug       bool
	settings    []SyncSetting
}

func NewConfig(c *cli.Context, debug bool) *Config {
	port := c.GlobalInt("port")
	dataDir := c.GlobalString("data")
	gitDataDir := dataDir + "/" + "git"
	confDir := dataDir + "/" + "conf"
	indexedDir := dataDir + "/" + "indexed"

	indexerType := c.GlobalString("indexer")

	sizeLimit := c.Int64("sizeLimit")

	schedule := c.GlobalString("schedule")

	config := &Config{
		DataDir:     dataDir,
		GitDataDir:  gitDataDir,
		ConfDir:     confDir,
		IndexedDir:  indexedDir,
		Port:        port,
		IndexerType: indexerType,
		SizeLimit:   sizeLimit,
		Schedule:    schedule,
		Debug:       false,
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
	if err := os.MkdirAll(c.IndexedDir, 0644); err != nil {
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

	c.SyncAllSCM()
	c.reloadSettings()

	return nil
}

func (c *Config) SyncAllSCM() error {
	for i := range c.settings {
		setting := c.settings[i]
		err := setting.SyncSCM()

		if err == nil {
			// write!
			c.writeSetting(setting.GetName())
		}
	}
	return nil
}

func (c *Config) SyncSCM(organization string) error {
	setting, ok := c.findSyncSetting(organization)
	if ok {
		err := setting.SyncSCM()
		if err == nil {
			// write!
			c.writeSetting(setting.GetName())
		}
		return nil
	}
	return errors.Errorf(`Not found setting for "%s"`, organization)
}

func (c *Config) reloadSettings() error {
	list := []SyncSetting{}

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

		var setting SyncSetting
		if organizationSetting.Scm["type"] == "bitbucket" {
			setting = NewBitbucketOrganizationSetting(organizationSetting)
		} else {
			setting = SyncSetting(&organizationSetting)
		}

		list = append(list, setting)
	}

	// cache update
	c.settings = list

	return nil
}

func (c *Config) GetSettings() []SyncSetting {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	return c.settings
}

func (c *Config) FindSetting(organization string) (SyncSetting, bool) {
	settings := c.GetSettings()
	for _, setting := range settings {
		if setting.GetName() == organization {
			return setting, true
		}
	}
	return nil, false
}

type SyncSetting interface {
	GetName() string
	GetProjects() []ProjectSetting
	GetSCM() map[string]string
	SyncSCM() error
	AddRepository(project string, repositoryUrl string) error
	FindProjectSetting(project string) (*ProjectSetting, bool)
	FindRepositorySetting(project string, repository string) (*RepositorySetting, bool)
	JSON() (string, error)
}

type OrganizationSetting struct {
	Name     string            `json:"name"`
	Projects []ProjectSetting  `json:"projects,omitempty"`
	Scm      map[string]string `json:"scm,omitempty"`
}

func (o *OrganizationSetting) GetName() string {
	return o.Name
}

func (o *OrganizationSetting) GetProjects() []ProjectSetting {
	return o.Projects
}

func (o *OrganizationSetting) GetSCM() map[string]string {
	return o.Scm
}

func (o *OrganizationSetting) SyncSCM() error {
	// do nothing
	return nil
}

func (o *OrganizationSetting) AddRepository(project string, url string) error {
	projectSetting, ok := o.FindProjectSetting(project)
	if !ok {
		projectSetting = &ProjectSetting{
			Name: project,
			Repositories: []RepositorySetting{
				RepositorySetting{
					Url: url,
				},
			},
		}
		o.Projects = append(o.Projects, *projectSetting)
		return nil
	}

	repositorySetting, ok := o.FindRepositorySetting(project, repoUrlToName(url))
	if !ok {
		repositorySetting = &RepositorySetting{
			Url: url,
		}
		projectSetting.Repositories = append(projectSetting.Repositories, *repositorySetting)
		return nil

	} else {
		return errors.Errorf("The repository already exists %s:%s/%s", o.Name, project, repoUrlToName(url))
	}
}

func (o *OrganizationSetting) FindProjectSetting(project string) (*ProjectSetting, bool) {
	for i := range o.Projects {
		if o.Projects[i].Name == project {
			return &o.Projects[i], true
		}
	}
	return nil, false
}

func (o *OrganizationSetting) FindRepositorySetting(project string, repository string) (*RepositorySetting, bool) {
	projectSetting, ok := o.FindProjectSetting(project)
	if ok {
		for i := range projectSetting.Repositories {
			if projectSetting.Repositories[i].GetName() == repository {
				return &projectSetting.Repositories[i], true
			}
		}
	}
	return nil, false
}

func (o *OrganizationSetting) JSON() (string, error) {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), err
}

type ProjectSetting struct {
	Name         string              `json:"name"`
	Repositories []RepositorySetting `json:"repositories"`
}

type RepositorySetting struct {
	Url  string `json:"url"`
	name string `json:"-"`
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

func (c *Config) findSyncSetting(organization string) (SyncSetting, bool) {
	for i := range c.settings {
		if c.settings[i].GetName() == organization {
			return c.settings[i], true
		}
	}
	return nil, false
}

func (c *Config) GetIndexed(organization string, project string, repository string) Indexed {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	return c.readIndexed(organization, project, repository)
}

func (c *Config) readIndexed(organization string, project string, repository string) Indexed {
	fileName := c.getIndexedFilePath(organization, project, repository)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return Indexed{Organization: organization, Project: project, Repository: repository, Branches: make(BrancheIndexedMap), Tags: make(TagIndexedMap)}
	}
	var indexed Indexed
	json.Unmarshal(content, &indexed)

	return indexed
}

func (c *Config) AddSetting(organization string, scmOptions map[string]string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	setting, ok := c.findSyncSetting(organization)
	if ok {
		return errors.Errorf(`The "%s" setting already exists`, organization)
	}

	setting = &OrganizationSetting{
		Name: organization,
		Scm:  scmOptions,
	}
	c.settings = append(c.settings, setting)

	// write!
	c.writeSetting(organization)

	return nil
}

func (c *Config) AddRepositorySetting(organization string, project string, url string, scmOptions map[string]string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	setting, ok := c.findSyncSetting(organization)
	if !ok {
		setting = &OrganizationSetting{
			Name: organization,
			Scm:  scmOptions,
			Projects: []ProjectSetting{
				ProjectSetting{
					Name: project,
					Repositories: []RepositorySetting{
						RepositorySetting{
							Url: url,
						},
					},
				},
			},
		}
		c.settings = append(c.settings, setting)

	} else {
		err := setting.AddRepository(project, url)
		if err != nil {
			return err
		}
	}

	// write!
	c.writeSetting(organization)

	return nil
}

func (c *Config) UpdateIndexed(indexed Indexed) error {
	indexedFileMutex.Lock()
	defer indexedFileMutex.Unlock()

	err := c.writeIndexed(indexed)
	return err
}

func (c *Config) DeleteIndexed(organization string, project string, repository string, removeBranches []string, removeTags []string) error {
	indexedFileMutex.Lock()
	defer indexedFileMutex.Unlock()

	currentIndexed := c.readIndexed(organization, project, repository)

	for _, removeBranch := range removeBranches {
		delete(currentIndexed.Branches, removeBranch)
	}

	for _, removeTag := range removeTags {
		delete(currentIndexed.Tags, removeTag)
	}

	return c.writeIndexed(currentIndexed)
}

func (c *Config) writeSetting(organization string) error {
	setting, ok := c.findSyncSetting(organization)
	if ok {
		content, _ := json.MarshalIndent(setting, "", "  ")
		fileName := fmt.Sprintf("%s/%s.json", c.ConfDir, organization)
		return ioutil.WriteFile(fileName, content, os.ModePerm)
	} else {
		return errors.Errorf("Not found organization: %s", organization)
	}
}

type Indexed struct {
	Organization string            `json:"organization"`
	Project      string            `json:"project"`
	Repository   string            `json:"repository"`
	Branches     BrancheIndexedMap `json:"branches"`
	Tags         TagIndexedMap     `json:"tags"`
}

type BrancheIndexedMap map[string]string
type TagIndexedMap map[string]string

func (c *Config) writeIndexed(indexed Indexed) error {
	content, _ := json.MarshalIndent(indexed, "", "  ")
	fileName := c.getIndexedFilePath(indexed.Organization, indexed.Project, indexed.Repository)

	if err := os.MkdirAll(fmt.Sprintf("%s/%s/%s", c.IndexedDir, indexed.Organization, indexed.Project), 0644); err != nil {
		log.Fatalln(err)
	}

	if c.Debug {
		log.Println("Write indexd file.", fileName)
	}

	return ioutil.WriteFile(fileName, content, os.ModePerm)
}

func (c *Config) getIndexedFilePath(organization string, project string, repository string) string {
	return fmt.Sprintf("%s/%s/%s/%s.json", c.IndexedDir, organization, project, repository)
}
