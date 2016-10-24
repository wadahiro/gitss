package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

type Config struct {
	DataDir     string
	GitDataDir  string
	ConfDir     string
	Port        int
	IndexerType string
	SizeLimit   int64
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

	config := Config{
		DataDir:     dataDir,
		GitDataDir:  gitDataDir,
		ConfDir:     confDir,
		Port:        port,
		IndexerType: indexerType,
		SizeLimit:   sizeLimit,
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

type IndexConf struct {
	Url  string `json:"url"`
	Refs []Ref  `json:"refs"`
}

type Ref struct {
	Name   string `json:"name"`
	Latest string `json:"latest"`
}

func (c *Config) GetIndexedCommitID(latestIndex LatestIndex) (string, bool) {
	fileName := c.getFileName(latestIndex)
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
	fileName := c.getFileName(latestIndex)
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

func (c *Config) writeFile(latestIndex LatestIndex, indexConf IndexConf) error {
	fileName := c.getFileName(latestIndex)
	content, _ := json.MarshalIndent(indexConf, "", "  ")
	return ioutil.WriteFile(fileName, content, os.ModePerm)
}

func (c *Config) getDir(latestIndex LatestIndex) string {
	dir := fmt.Sprintf("%s/%s/%s", c.ConfDir, latestIndex.Organization, latestIndex.Project)
	return dir
}

func (c *Config) getFileName(latestIndex LatestIndex) string {
	fileName := fmt.Sprintf("%s/%s/%s/%s.json", c.ConfDir, latestIndex.Organization, latestIndex.Project, latestIndex.Repository)
	return fileName
}
