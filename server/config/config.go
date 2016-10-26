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

	"github.com/pkg/errors"

	"github.com/codegangsta/cli"
)

type Config struct {
	DataDir     string
	GitDataDir  string
	ConfDir     string
	Port        int
	IndexerType string
	SizeLimit   int64
	Schedule    string
	Debug       bool
}

type LatestIndex struct {
	Organization string `json:"organization"`
	Project      string `json:"project"`
	Repository   string `json:"repository"`
	Ref          string `json:"ref"`
}

func NewConfig(c *cli.Context, debug bool) Config {
	port := c.GlobalInt("port")
	dataDir := c.GlobalString("data")
	gitDataDir := dataDir + "/" + "git"
	confDir := dataDir + "/" + "conf"

	indexerType := c.GlobalString("indexer")

	sizeLimit := c.Int64("sizeLimit")

	schedule := c.GlobalString("schedule")

	config := Config{
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
}

func (c *Config) GetAllIndexConf() ([]IndexConfWrapper, error) {
	list := []IndexConfWrapper{}

	err := filepath.Walk(c.ConfDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		log.Println(path)
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			paths := strings.Split(path, string(os.PathSeparator))
			if len(paths) < 3 {
				return nil
			}

			organization := paths[len(paths)-3]
			project := paths[len(paths)-2]
			repository := strings.Split(info.Name(), ".json")[0]

			b, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Errorf("Not found config: %s", path) // NotFound
			}
			var indexConf IndexConfWrapper
			json.Unmarshal(b, &indexConf)

			indexConf.Organization = organization
			indexConf.Project = project
			indexConf.Repository = repository

			list = append(list, indexConf)
		}
		return nil
	})

	return list, err
}

type IndexConfWrapper struct {
	IndexConf
	Organization string
	Project      string
	Repository   string
}

type IndexConf struct {
	Url  string `json:"url"`
	Refs []Ref  `json:"refs"`
}

type Ref struct {
	Name   string `json:"name"`
	Latest string `json:"latest"`
}

func (c *Config) GetRefs(organization string, project string, repository string) ([]Ref, error) {
	fileName := c.getFileName(organization, project, repository)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Errorf("Not found config: %s", fileName) // NotFound
	}
	var indexConf IndexConf
	json.Unmarshal(b, &indexConf)

	return indexConf.Refs, nil
}

func (c *Config) GetIndexedCommitID(latestIndex LatestIndex) (string, bool) {
	fileName := c.getFileNameByLatestIndex(latestIndex)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", true // NotFound
	}
	var indexConf IndexConf
	json.Unmarshal(b, &indexConf)

	for i := range indexConf.Refs {
		ref := indexConf.Refs[i]
		if ref.Name == latestIndex.Ref {
			return ref.Latest, false
		}
	}

	return "", true
}

func (c *Config) UpdateLatestIndex(url string, latestIndex LatestIndex, commitId string) error {
	fileName := c.getFileNameByLatestIndex(latestIndex)
	b, err := ioutil.ReadFile(fileName)

	var indexConf IndexConf
	if err != nil {
		// Create
		indexConf = IndexConf{
			Url: url,
			Refs: []Ref{
				Ref{
					Name:   latestIndex.Ref,
					Latest: commitId,
				},
			},
		}
		os.MkdirAll(c.getDir(latestIndex), 0644)
		return c.writeFile(latestIndex, indexConf)

	} else {
		// Update latest
		json.Unmarshal(b, &indexConf)

		for i := range indexConf.Refs {
			ref := indexConf.Refs[i]
			if ref.Name == latestIndex.Ref {
				indexConf.Refs[i].Latest = commitId

				return c.writeFile(latestIndex, indexConf)
			}
		}
		// Add ref
		indexConf.Refs = append(indexConf.Refs, Ref{
			Name:   latestIndex.Ref,
			Latest: commitId,
		})
		return c.writeFile(latestIndex, indexConf)
	}

	return nil
}

func (c *Config) DeleteLatestIndexRefs(organization string, project string, repository string, refs []string) error {
	fileName := c.getFileName(organization, project, repository)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	var indexConf IndexConf
	// Update latest
	json.Unmarshal(b, &indexConf)

	newRefs := []Ref{}
	for i := range indexConf.Refs {
		ref := indexConf.Refs[i]

		found := false
		for _, removeRef := range refs {
			if ref.Name == removeRef {
				found = true
				break
			}
		}
		if !found {
			newRefs = append(newRefs, indexConf.Refs[i])
		}
	}

	// Update ref
	indexConf.Refs = newRefs
	return c.writeFileByFineName(fileName, indexConf)
}

func (c *Config) writeFile(latestIndex LatestIndex, indexConf IndexConf) error {
	fileName := c.getFileNameByLatestIndex(latestIndex)
	return c.writeFileByFineName(fileName, indexConf)
}

func (c *Config) writeFileByFineName(fileName string, indexConf IndexConf) error {
	content, _ := json.MarshalIndent(indexConf, "", "  ")
	return ioutil.WriteFile(fileName, content, os.ModePerm)
}

func (c *Config) getDir(latestIndex LatestIndex) string {
	dir := fmt.Sprintf("%s/%s/%s", c.ConfDir, latestIndex.Organization, latestIndex.Project)
	return dir
}

func (c *Config) getFileNameByLatestIndex(latestIndex LatestIndex) string {
	return c.getFileName(latestIndex.Organization, latestIndex.Project, latestIndex.Repository)
}

func (c *Config) getFileName(organization string, project string, repository string) string {
	fileName := fmt.Sprintf("%s/%s/%s/%s.json", c.ConfDir, organization, project, repository)
	return fileName
}
