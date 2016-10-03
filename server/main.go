package main

import (
	"log"
	"os"

	"strconv"

	"github.com/codegangsta/cli"

	"github.com/wadahiro/gits/server/importer"
	"github.com/wadahiro/gits/server/indexer"
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
				Usage: "data directory",
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

	log.Println("-------------- GitS Server --------------")
	log.Println("VERSION: ", Version)
	log.Println("COMMIT_HASH: ", CommitHash)
	log.Println("DATA_DIR: ", dataDir)
	log.Println("PORT: ", port)
	log.Println("DEBUG_MODE: ", debugMode)
	log.Println("-----------------------------------------")

	if err := os.MkdirAll(dataDir, 0644); err != nil {
		log.Fatalln(err)
	}

	// service.RunSyncScheduler(repo)
	
	indexer := indexer.NewESIndexer()

	initRouter(indexer, port, debugMode, dataDir)

	log.Println("Started GitS Server.")
}

func ImportGitRepository(c *cli.Context) {
	debugMode := isDebugMode()

	dataDir := c.GlobalString("data")

	if len(c.Args()) != 2 {
		log.Fatalln("Please specified [project name] [git repository url]")
	}

	projectName := c.Args()[0]
	gitRepoUrl := c.Args()[1]

	log.Println("-------------- GitS Import Git Repository --------------")
	log.Println("VERSION: ", Version)
	log.Println("COMMIT_HASH: ", CommitHash)
	log.Println("DATA_DIR: ", dataDir)
	log.Println("DEBUG_MODE: ", debugMode)
	log.Println("PROJECT_NAME: ", projectName)
	log.Println("GIT_REPOSITORY_URL: ", gitRepoUrl)
	log.Println("--------------------------------------------------------")

	indexer := indexer.NewESIndexer()
	importer := importer.NewGitImporter(dataDir, indexer, debugMode)
	importer.Run(projectName, gitRepoUrl)
}

func isDebugMode() bool {
	return BuildTarget == "develop"
}
