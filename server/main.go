package main

import (
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/codegangsta/cli"

	"github.com/wadahiro/gitss/server/config"
	"github.com/wadahiro/gitss/server/importer"
	"github.com/wadahiro/gitss/server/indexer"
	"github.com/wadahiro/gitss/server/repo"
	"github.com/wadahiro/gitss/server/service"
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
	app.Usage = "Git Source Search."
	app.Version = Version
	app.Author = "Hiroyuki Wada"
	app.Email = "wadahiro@gmail.com"
	app.Commands = []cli.Command{
		{
			Name:   "server",
			Usage:  "Run GitSS",
			Action: RunServer,
			Flags: []cli.Flag{
				cli.Int64Flag{
					Name:  "sizeLimit",
					Value: 1024 * 1024, //1MB
					Usage: "Indexing limit file size",
				},
			},
		},
		{
			Name:      "sync",
			Usage:     "Sync all git repositories",
			ArgsUsage: "",
			Action:    SyncAll,
			Flags: []cli.Flag{
				cli.Int64Flag{
					Name:  "sizeLimit",
					Value: 1024 * 1024, //1MB
					Usage: "Indexing limit file size",
				},
			},
		},
		{
			Name:      "add",
			Usage:     "Add a git repository",
			ArgsUsage: "",
			Action:    AddGitRepository,
		},
		{
			Name:      "bitbucket",
			Usage:     "Bitbucket server related commands",
			ArgsUsage: "",
			Subcommands: []cli.Command{
				{
					Name:      "add",
					Usage:     "Add a sync setting with the bitbucket server",
					ArgsUsage: "[organization name] [bitbucket server url] --user=[username] --password=[password] --include-projects=<...> --exclude-projects=<...> --include-repositories=<...> --exclude-repositories=<...> --include-branches=<...> --exclude-branches=<...> --include-tags=<...> --exclude-tags=<...>",
					Action:    AddBitbucketServerSetting,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user",
							Value: "",
							Usage: "Set username for the bitbucket server if you'd like to use Basic Authentication",
						},
						cli.StringFlag{
							Name:  "password",
							Value: "",
							Usage: "Set password for the bitbucket server if you'd like to use Basic Authentication",
						},
						cli.StringFlag{
							Name:  "include-projects",
							Value: ".*",
							Usage: "Set regex pattern of the name of the projects which you'd like to include",
						},
						cli.StringFlag{
							Name:  "exclude-projects",
							Value: ".*",
							Usage: "Set regex pattern of the name of the projects which you'd like to exclude",
						},
						cli.StringFlag{
							Name:  "include-repositories",
							Value: ".*",
							Usage: "Set regex pattern of the name of the repositories which you'd like to include",
						},
						cli.StringFlag{
							Name:  "exclude-repositories",
							Value: ".*",
							Usage: "Set regex pattern of the name of the repositories which you'd like to exclude",
						},
						cli.StringFlag{
							Name:  "include-branches",
							Value: ".*",
							Usage: "Set regex pattern of the name of the branches which you'd like to include",
						},
						cli.StringFlag{
							Name:  "exclude-branches",
							Value: ".*",
							Usage: "Set regex pattern of the name of the branches which you'd like to exclude",
						},
						cli.StringFlag{
							Name:  "include-tags",
							Value: ".*",
							Usage: "Set regex pattern of the name of the tags which you'd like to include",
						},
						cli.StringFlag{
							Name:  "exclude-tags",
							Value: ".*",
							Usage: "Set regex pattern of the name of the tags which you'd like to exclude",
						},
					},
				},
			},
		},
		{
			Name:      "import",
			Usage:     "Import a git repository",
			ArgsUsage: "[organization name] [project name] [git repository url]",
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
		cli.StringFlag{
			Name:  "schedule",
			Value: "0 */10 * * * *",
			Usage: "Sync schedule",
		},
	}
	app.Run(args)
}

func RunServer(c *cli.Context) {
	debugMode := isDebugMode()

	config := config.NewConfig(c, debugMode)

	log.Println("-------------- GitSS --------------------")
	log.Println("VERSION: ", Version)
	log.Println("COMMIT_HASH: ", CommitHash)
	log.Println("DATA_DIR: ", config.DataDir)
	log.Println("INDEXER_TYPE: ", config.IndexerType)
	log.Println("PORT: ", strconv.Itoa(config.Port))
	log.Println("SCHEDULE: ", config.Schedule)
	log.Println("DEBUG_MODE: ", debugMode)
	log.Println("-----------------------------------------")

	reader := repo.NewGitRepoReader(config)
	indexer := newIndexer(config, reader)
	importer := importer.NewGitImporter(config, indexer)
	service.RunSyncScheduler(config, importer)

	initRouter(config, indexer)
}

func SyncAll(c *cli.Context) {
	debugMode := isDebugMode()

	config := config.NewConfig(c, debugMode)
	reader := repo.NewGitRepoReader(config)
	indexer := newIndexer(config, reader)
	importer := importer.NewGitImporter(config, indexer)

	service.RunSync(config, importer)
}

func AddGitRepository(c *cli.Context) {
	debugMode := isDebugMode()

	if len(c.Args()) != 3 {
		log.Fatalln("Please specified [organization name] [project name] [git repository url]")
	}

	config := config.NewConfig(c, debugMode)

	organization := c.Args()[0]
	projectName := c.Args()[1]
	gitRepoUrl := c.Args()[2]

	err := config.AddRepositorySetting(organization, projectName, gitRepoUrl, nil)
	if err != nil {
		log.Println(err)
	}
}

func AddBitbucketServerSetting(c *cli.Context) {
	debugMode := isDebugMode()

	if len(c.Args()) != 2 {
		log.Fatalln("Please specified [organization name] [bitbucket server url]")
	}

	config := config.NewConfig(c, debugMode)

	organization := c.Args()[0]
	bitbucketUrl := c.Args()[1]
	user := c.String("user")
	password := c.String("password")
	includeProjects := regex(c.String("include-projects"))
	excludeProjects := regex(c.String("exclude-projects"))
	includeRepositories := regex(c.String("include-repositories"))
	excludeRepositories := regex(c.String("exclude-repositories"))
	includeBranches := regex(c.String("include-branches"))
	excludeBranches := regex(c.String("exclude-branches"))
	includeTags := regex(c.String("include-tags"))
	excludeTags := regex(c.String("exclude-tags"))

	scmOptions := make(map[string]string)
	scmOptions["type"] = "bitbucket"
	scmOptions["url"] = bitbucketUrl
	scmOptions["user"] = user
	scmOptions["password"] = password
	scmOptions["includeProjects"] = includeProjects
	scmOptions["excludeProjects"] = excludeProjects
	scmOptions["includeRepositories"] = includeRepositories
	scmOptions["excludeRepositories"] = excludeRepositories

	err := config.AddSetting(organization, scmOptions, includeBranches, excludeBranches, includeTags, excludeTags)
	if err != nil {
		log.Println(err)
	}
}

func regex(pattern string) string {
	regexp.MustCompile(pattern)
	return pattern
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

	log.Println("-------------- GitSS Import Git Repository -------------")
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

	config.AddRepositorySetting(organization, projectName, gitRepoUrl, nil)

	reader := repo.NewGitRepoReader(config)
	indexer := newIndexer(config, reader)
	importer := importer.NewGitImporter(config, indexer)
	importer.Run(organization, projectName, gitRepoUrl)
}

func isDebugMode() bool {
	return BuildTarget == "develop"
}

func newIndexer(config *config.Config, reader *repo.GitRepoReader) indexer.Indexer {
	switch config.IndexerType {
	case "bleve":
		return indexer.NewBleveIndexer(config, reader)
	case "es":
		return indexer.NewESIndexer(config, reader)
	}
	panic("Unknown Indexer type")
}
