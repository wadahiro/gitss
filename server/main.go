package main

import (
	"log"
	"os"
	"strconv"

	"github.com/codegangsta/cli"

	"github.com/wadahiro/gitss/server/config"
	"github.com/wadahiro/gitss/server/importer"
	"github.com/wadahiro/gitss/server/indexer"
	"github.com/wadahiro/gitss/server/repo"
)

var CommitHash = "WORKSPACE" // inject by LDFLAGS build option
var Version = "SNAPSHOT"     // inject by LDFLAGS build option
var BuildTarget = "develop"  // inject by LDFLAGS build option

func main() {
	args := os.Args
	if BuildTarget == "develop" && len(args) == 1 {
		args = []string{"gitss", "server"}
	}

	app := cli.NewApp()
	app.Name = "GitSS"
	app.Usage = "Code Search Server for Git repositories."
	app.Version = Version
	app.Author = "Hiroyuki Wada"
	app.Email = "wadahiro@gmail.com"
	app.Commands = []cli.Command{
		{
			Name:   "server",
			Usage:  "Run GitSS",
			Action: RunServer,
		},
		{
			Name:      "import",
			Usage:     "Import a Git repository",
			ArgsUsage: "[project name] [git repository url]",
			Action:    ImportGitRepository,
			Flags: []cli.Flag{
				cli.Int64Flag{
					Name:  "sizeLimit",
					Value: 1024 * 1024, //1MB
					Usage: "Indexing limit file size",
				},
			},
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
	app.Run(args)
}

func RunServer(c *cli.Context) {
	debugMode := isDebugMode()

	config := config.NewConfig(c, debugMode)

	log.Println("-------------- GitS Server --------------")
	log.Println("VERSION: ", Version)
	log.Println("COMMIT_HASH: ", CommitHash)
	log.Println("DATA_DIR: ", config.DataDir)
	log.Println("INDEXER_TYPE: ", config.IndexerType)
	log.Println("PORT: ", strconv.Itoa(config.Port))
	log.Println("DEBUG_MODE: ", debugMode)
	log.Println("-----------------------------------------")

	// service.RunSyncScheduler(repo)

	reader := repo.NewGitRepoReader(config)
	indexer := newIndexer(config, reader)
	initRouter(config, indexer)

	log.Println("Started GitS Server.")
}

func ImportGitRepository(c *cli.Context) {
	debugMode := isDebugMode()

	if len(c.Args()) != 3 {
		log.Fatalln("Please specified [organization name] [project name] [git repository url]")
	}

	config := config.NewConfig(c, debugMode)

	organization := c.Args()[0]
	projectName := c.Args()[1]
	gitRepoUrl := c.Args()[2]

	log.Println("-------------- GitS Import Git Repository --------------")
	log.Println("VERSION: ", Version)
	log.Println("COMMIT_HASH: ", CommitHash)
	log.Println("DATA_DIR: ", config.DataDir)
	log.Println("INDEXER_TYPE: ", config.IndexerType)
	log.Println("DEBUG_MODE: ", config.Debug)
	log.Println("ORGANIZATION_NAME: ", organization)
	log.Println("PROJECT_NAME: ", projectName)
	log.Println("GIT_REPOSITORY_URL: ", gitRepoUrl)
	log.Println("SIZE_LIMIT: ", config.SizeLimit)
	log.Println("--------------------------------------------------------")

	reader := repo.NewGitRepoReader(config)
	indexer := newIndexer(config, reader)
	importer := importer.NewGitImporter(config, indexer)
	importer.Run(organization, projectName, gitRepoUrl)
}

func isDebugMode() bool {
	return BuildTarget == "develop"
}

func newIndexer(config config.Config, reader *repo.GitRepoReader) indexer.Indexer {
	switch config.IndexerType {
	case "bleve":
		return indexer.NewBleveIndexer(config, reader)
	case "es":
		return indexer.NewESIndexer(config, reader)
	}
	panic("Unknown Indexer type")
}
