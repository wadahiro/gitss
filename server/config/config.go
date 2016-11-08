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
	settings    []SyncSetting
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

type SyncSetting interface {
	GetName() string
	GetProjects() []ProjectSetting
	GetSCM() map[string]string
	SyncSCM() error
	AddRepository(project string, repositoryUrl string) error
	DeleteRefs(project string, repository string, removeRefs []string)
	FindProjectSetting(project string) (*ProjectSetting, bool)
	FindRepositorySetting(project string, repository string) (*RepositorySetting, bool)
	FindRefSetting(project string, repository string, refs string) (*RefSetting, bool)
	JSON() (string, error)
}

type OrganizationSetting struct {
	Name     string            `json:"name"`
	Projects []ProjectSetting  `json:"projects"`
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
					Url:  url,
					Refs: []RefSetting{},
				},
			},
		}
		o.Projects = append(o.Projects, *projectSetting)
		return nil
	}

	repositorySetting, ok := o.FindRepositorySetting(project, repoUrlToName(url))
	if !ok {
		repositorySetting = &RepositorySetting{
			Url:  url,
			Refs: []RefSetting{},
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

func (o *OrganizationSetting) FindRefSetting(project string, repository string, ref string) (*RefSetting, bool) {
	repositorySetting, ok := o.FindRepositorySetting(project, repository)
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

func (o *OrganizationSetting) DeleteRefs(project string, repository string, removeRefs []string) {
	repositorySetting, ok := o.FindRepositorySetting(project, repository)
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
	}
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
	Url  string       `json:"url"`
	name string       `json:"-"`
	Refs []RefSetting `json:"refs,omitempty"`
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

func (c *Config) findSyncSetting(organization string) (SyncSetting, bool) {
	for i := range c.settings {
		if c.settings[i].GetName() == organization {
			return c.settings[i], true
		}
	}
	return nil, false
}

func (c *Config) GetRefs(organization string, project string, repository string) ([]RefSetting, bool) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	syncSetting, ok := c.findSyncSetting(organization)
	if ok {
		repositorySetting, ok := syncSetting.FindRepositorySetting(project, repository)
		if ok {
			return repositorySetting.Refs, true
		}
	}

	return nil, false
}

func (c *Config) GetIndexedCommitID(organization string, project string, repository string, ref string) (string, bool) {
	syncSetting, ok := c.findSyncSetting(organization)
	if ok {
		refSetting, ok := syncSetting.FindRefSetting(project, repository, ref)
		if ok {
			return refSetting.Latest, true
		}
	}
	return "", false
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
							Url:  url,
							Refs: []RefSetting{},
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

func (c *Config) UpdateLatestIndex(url string, organization string, project string, repository string, ref string, commitId string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// update case
	syncSetting, ok := c.findSyncSetting(organization)
	if ok {
		refSetting, ok := syncSetting.FindRefSetting(project, repository, ref)
		if ok {
			refSetting.Latest = commitId
		} else {
			// add ref case
			repositorySetting, ok := syncSetting.FindRepositorySetting(project, repository)
			if ok {
				repositorySetting.Refs = append(repositorySetting.Refs, RefSetting{Name: ref, Latest: commitId})
			}
		}
		// write
		return c.writeSetting(organization)
	}
	return nil
}

func (c *Config) DeleteLatestIndexRefs(organization string, project string, repository string, removeRefs []string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	syncSetting, ok := c.findSyncSetting(organization)
	if ok {
		syncSetting.DeleteRefs(project, repository, removeRefs)
		return c.writeSetting(organization)
	}

	return errors.Errorf("Not found repository setting for %s:%s/%s", organization, project, repository)
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
