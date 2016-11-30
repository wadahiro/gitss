package service

import (
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron"
	// "github.com/wadahiro/gitss/server/model"
	"github.com/wadahiro/gitss/server/config"
	"github.com/wadahiro/gitss/server/importer"
	"github.com/wadahiro/gitss/server/util"
)

var mutex = new(sync.Mutex)
var scheduler *cron.Cron

func RunSyncScheduler(config *config.Config, importer *importer.GitImporter) {
	mutex.Lock()
	defer mutex.Unlock()

	if scheduler != nil {
		log.Println("Stop sync schduler.")
		scheduler.Stop()
	}
	scheduler = cron.New()

	spec := config.Schedule

	log.Printf("Setup sync job. spec: %s\n", spec)

	job := func() {
		log.Println("Start sync job.")
		defer log.Println("End Start sync job.")

		RunSyncAll(config, importer)
	}

	scheduler.AddFunc(spec, job)

	scheduler.Start()
	fmt.Println("Started sync schduler.")
}

func RunSync(config *config.Config, importer *importer.GitImporter, organization, project, repository string) {
	config.Sync()

	setting, ok := config.FindSetting(organization)
	if !ok {
		log.Printf("Not found organization: %s\n", organization)
		return
	}

	projectSetting, ok := setting.FindProjectSetting(project)
	if !ok {
		log.Printf("Not found project: %v\n", project)
		return
	}

	repositorySetting, ok := setting.FindRepositorySetting(project, repository)
	if !ok {
		log.Printf("Not found repository: %s\n", repository)
		return
	}

	log.Printf("Sync for %s:%s/%s\n", setting.GetName(), projectSetting.Name, repositorySetting.GetName())

	importer.Run(setting.GetName(), projectSetting.Name, repositorySetting.Url)
}

func RunSyncAll(config *config.Config, importer *importer.GitImporter) {
	config.Sync()

	settings := config.GetSettings()

	workers := util.GenWorkers(2)
	var wg sync.WaitGroup

	for i := range settings {
		setting := settings[i]

		for _, project := range setting.GetProjects() {

			for k := range project.Repositories {
				repository := project.Repositories[k]

				log.Printf("Sync for %s:%s/%s\n", setting.GetName(), project.Name, repository.GetName())

				wg.Add(1)
				workers <- func() {
					defer wg.Done()
					importer.Run(setting.GetName(), project.Name, repository.Url)
				}
			}
		}
	}
	wg.Wait()
}
