package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

type Config struct {
	DataDir     string
	GitDataDir  string
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

	indexerType := c.GlobalString("indexer")

	sizeLimit := c.Int64("sizeLimit")

	config := Config{
		DataDir:     dataDir,
		GitDataDir:  gitDataDir,
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

func (c *Config) GetIndexedCommitID(latestIndex LatestIndex) (string, bool) {
	fileName := c.getFileName(latestIndex)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", true // NotFound
	}
	return string(b), false
}

func (c *Config) UpdateLatestIndex(latestIndex LatestIndex, commitId string) error {
	fileName := c.getFileName(latestIndex)
	err := ioutil.WriteFile(fileName, []byte(commitId), os.ModePerm)
	return err
}

func (c *Config) getFileName(latestIndex LatestIndex) string {
	fileName := fmt.Sprintf("%s/%s/%s/%s.%s.indexed", c.GitDataDir, latestIndex.Organization, latestIndex.Project, latestIndex.Repository, latestIndex.Ref)
	return fileName
}
