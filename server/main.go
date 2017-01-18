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
				cli.IntFlag{
					Name:  "port",
					Value: 3000,
					Usage: "port number",
				},
				cli.StringFlag{
					Name:  "schedule",
					Value: "0 */10 * * * *",
					Usage: "Sync schedule",
				},
			},
		},
		{
			Name:      "sync",
			Usage:     "Sync git repository",
			ArgsUsage: "[ORGANIZATION] [PROJECT] [REPOSITORY]",
			Action:    Sync,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all",
					Usage: "Sync all git repositories",
				},
				cli.BoolFlag{
					Name:  "skip-git-sync",
					Usage: "Skip syncing with the remote git repository",
				},
				cli.BoolFlag{
					Name:  "skip-index",
					Usage: "Skip indexing the repository",
				},
			},
		},
		{
			Name:      "add",
			Usage:     "Add a sync setting",
			ArgsUsage: "ORGANIZATION PROJECT GIT_REPOSITORY_URL",
			Action:    AddGitRepository,
			Flags: []cli.Flag{
				cli.Int64Flag{
					Name:  "sizeLimit",
					Value: 1024 * 1024, //1MB
					Usage: "Indexing limit file size (byte)",
				},
				cli.StringFlag{
					Name:  "include-branches",
					Value: ".*",
					Usage: "Set regex pattern of the name of the branches which you'd like to include",
				},
				cli.StringFlag{
					Name:  "exclude-branches",
					Value: "",
					Usage: "Set regex pattern of the name of the branches which you'd like to exclude",
				},
				cli.StringFlag{
					Name:  "include-tags",
					Value: ".*",
					Usage: "Set regex pattern of the name of the tags which you'd like to include",
				},
				cli.StringFlag{
					Name:  "exclude-tags",
					Value: "",
					Usage: "Set regex pattern of the name of the tags which you'd like to exclude",
				},
			},
		},
		{
			Name:      "bitbucket",
			Usage:     "Bitbucket server related commands",
			ArgsUsage: "",
			Subcommands: []cli.Command{
				{
					Name:      "add",
					Usage:     "Add a sync setting with the bitbucket server",
					ArgsUsage: "ORGANIZATION BITBUCKET_SERVER_URL",
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
						cli.Int64Flag{
							Name:  "sizeLimit",
							Value: 1024 * 1024, //1MB
							Usage: "Indexing limit file size (byte)",
						},
						cli.StringFlag{
							Name:  "include-projects",
							Value: ".*",
							Usage: "Set regex pattern of the name of the projects which you'd like to include",
						},
						cli.StringFlag{
							Name:  "exclude-projects",
							Value: "",
							Usage: "Set regex pattern of the name of the projects which you'd like to exclude",
						},
						cli.StringFlag{
							Name:  "include-repositories",
							Value: ".*",
							Usage: "Set regex pattern of the name of the repositories which you'd like to include",
						},
						cli.StringFlag{
							Name:  "exclude-repositories",
							Value: "",
							Usage: "Set regex pattern of the name of the repositories which you'd like to exclude",
						},
						cli.StringFlag{
							Name:  "include-branches",
							Value: ".*",
							Usage: "Set regex pattern of the name of the branches which you'd like to include",
						},
						cli.StringFlag{
							Name:  "exclude-branches",
							Value: "",
							Usage: "Set regex pattern of the name of the branches which you'd like to exclude",
						},
						cli.StringFlag{
							Name:  "include-tags",
							Value: ".*",
							Usage: "Set regex pattern of the name of the tags which you'd like to include",
						},
						cli.StringFlag{
							Name:  "exclude-tags",
							Value: "",
							Usage: "Set regex pattern of the name of the tags which you'd like to exclude",
						},
					},
				},
			},
		},
	}
	app.Flags = []cli.Flag{
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

func Sync(c *cli.Context) error {
	debugMode := isDebugMode()

	all := c.Bool("all")

	config := config.NewConfig(c, debugMode)
	reader := repo.NewGitRepoReader(config)
	indexer := newIndexer(config, reader)
	importer := importer.NewGitImporter(config, indexer)

	if all {
		service.RunSyncAll(config, importer)

	} else {
		if len(c.Args()) != 3 {
			return cli.NewExitError("Please specified "+c.Command.ArgsUsage, 1)
		}

		organization := c.Args()[0]
		project := c.Args()[1]
		repository := c.Args()[2]

		service.RunSync(config, importer, organization, project, repository)
	}
	return nil
}

func AddGitRepository(c *cli.Context) error {
	debugMode := isDebugMode()

	if len(c.Args()) != 3 {
		return cli.NewExitError("Please specified "+c.Command.ArgsUsage, 1)
	}

	config := config.NewConfig(c, debugMode)

	organization := c.Args()[0]
	projectName := c.Args()[1]
	gitRepoUrl := c.Args()[2]

	sizeLimit := c.Int64("sizeLimit")

	includeBranches := regex(c.String("include-branches"))
	excludeBranches := regex(c.String("exclude-branches"))
	includeTags := regex(c.String("include-tags"))
	excludeTags := regex(c.String("exclude-tags"))

	err := config.AddRepositorySetting(organization, projectName, gitRepoUrl, nil, sizeLimit, includeBranches, excludeBranches, includeTags, excludeTags)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func AddBitbucketServerSetting(c *cli.Context) error {
	debugMode := isDebugMode()

	if len(c.Args()) != 2 {
		return cli.NewExitError("Please specified "+c.Command.ArgsUsage, 1)
	}

	config := config.NewConfig(c, debugMode)

	organization := c.Args()[0]
	bitbucketUrl := c.Args()[1]

	user := c.String("user")
	password := c.String("password")
	sizeLimit := c.Int64("sizeLimit")
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

	err := config.AddSetting(organization, scmOptions, sizeLimit, includeBranches, excludeBranches, includeTags, excludeTags)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func regex(pattern string) string {
	regexp.MustCompile(pattern)
	return pattern
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
