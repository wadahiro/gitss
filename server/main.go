package main

import (
	"log"
	"os"

	"strconv"

	"github.com/codegangsta/cli"

	"github.com/wadahiro/gitss/server/importer"
	"github.com/wadahiro/gitss/server/indexer"
	"github.com/wadahiro/gitss/server/repo"
)

var CommitHash = "WORKSPACE" // inject by LDFLAGS build option
var Version = "SNAPSHOT"     // inject by LDFLAGS build option
var BuildTarget = "develop"  // inject by LDFLAGS build option

func main() {
	if BuildTarget == "develop" {
		RunServer(nil)
	} else {
		app := cli.NewApp()
		app.Name = "GitS"
		app.Usage = "Code Search for Git repositories."
		app.Version = Version
		app.Author = "Hiroyuki Wada"
		app.Email = "wadahiro@gmail.com"
		app.Commands = []cli.Command{
			{
				Name:   "server",
				Usage:  "Run GitS server",
				Action: RunServer,
			},
			{
				Name:      "import",
				Usage:     "Import a Git repository",
				ArgsUsage: "[project name] [git repository url]",
				Action:    ImportGitRepository,
			},
		}
		app.Flags = []cli.Flag{
			cli.IntFlag{
				Name:  "port",
				Value: 3000,
				Usage: "port number",
			},
			cli.StringFlag{
				Name:  "data",
				Value: "./data",
				Usage: "Data directory",
			},
			cli.StringFlag{
				Name:  "indexer",
				Value: "bleve",
				Usage: "Indexer implementation",
			},
		}
		app.Run(os.Args)
	}
}

func RunServer(c *cli.Context) {
	debugMode := isDebugMode()

	port := "3000"
	if c != nil && c.GlobalInt("port") != 0 {
		port = strconv.Itoa(c.GlobalInt("port"))
	}

	dataDir := "./data"
	if c != nil {
		dataDir = c.GlobalString("data")
	}
	gitDataDir := dataDir + "/" + "git"

	indexerType := "bleve"
	if c != nil {
		indexerType = c.GlobalString("indexer")
	}

	log.Println("-------------- GitS Server --------------")
	log.Println("VERSION: ", Version)
	log.Println("COMMIT_HASH: ", CommitHash)
	log.Println("DATA_DIR: ", dataDir)
	log.Println("INDEXER_TYPE: ", indexerType)
	log.Println("PORT: ", port)
	log.Println("DEBUG_MODE: ", debugMode)
	log.Println("-----------------------------------------")

	if err := os.MkdirAll(gitDataDir, 0644); err != nil {
		log.Fatalln(err)
	}

	// service.RunSyncScheduler(repo)

	reader := repo.NewGitRepoReader(gitDataDir, debugMode)
	indexer := newIndexer(indexerType, reader, dataDir, debugMode)

	initRouter(indexer, port, debugMode, gitDataDir)

	log.Println("Started GitS Server.")
}

func ImportGitRepository(c *cli.Context) {
	debugMode := isDebugMode()

	dataDir := c.GlobalString("data")
	gitDataDir := dataDir + "/" + "git"

	indexerType := "bleve"
	if c != nil {
		indexerType = c.GlobalString("indexer")
	}

	if len(c.Args()) != 3 {
		log.Fatalln("Please specified [organization name] [project name] [git repository url]")
	}

	organization := c.Args()[0]
	projectName := c.Args()[1]
	gitRepoUrl := c.Args()[2]

	log.Println("-------------- GitS Import Git Repository --------------")
	log.Println("VERSION: ", Version)
	log.Println("COMMIT_HASH: ", CommitHash)
	log.Println("DATA_DIR: ", dataDir)
	log.Println("INDEXER_TYPE: ", indexerType)
	log.Println("DEBUG_MODE: ", debugMode)
	log.Println("ORGANIZATION_NAME: ", organization)
	log.Println("PROJECT_NAME: ", projectName)
	log.Println("GIT_REPOSITORY_URL: ", gitRepoUrl)
	log.Println("--------------------------------------------------------")

	reader := repo.NewGitRepoReader(gitDataDir, debugMode)
	indexer := newIndexer(indexerType, reader, dataDir, debugMode)
	importer := importer.NewGitImporter(gitDataDir, indexer, debugMode)
	importer.Run(organization, projectName, gitRepoUrl)
}

func isDebugMode() bool {
	return BuildTarget == "develop"
}

func newIndexer(indexerType string, reader *repo.GitRepoReader, dataDir string, debugMode bool) indexer.Indexer {
	switch indexerType {
	case "bleve":
		return indexer.NewBleveIndexer(reader, dataDir+"/bleve_index", debugMode)
	case "es":
		return indexer.NewESIndexer(reader, debugMode)
	}
	panic("Unknown Indexer type")
}
