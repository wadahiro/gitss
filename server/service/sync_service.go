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

		RunSync(config, importer)
	}

	scheduler.AddFunc(spec, job)

	scheduler.Start()
	fmt.Println("Started sync schduler.")
}

func RunSync(config *config.Config, importer *importer.GitImporter) {
	config.Sync()

	organizations := config.GetSettings()

	workers := util.GenWorkers(2)
	var wg sync.WaitGroup

	for i := range organizations {
		organization := organizations[i]

		for j := range organization.Projects {
			project := organization.Projects[j]

			for k := range project.Repositories {
				repository := project.Repositories[k]

				log.Printf("Sync for %s:%s/%s\n", organization.Name, project.Name, repository.GetName())

				wg.Add(1)
				workers <- func() {
					defer wg.Done()
					importer.Run(organization.Name, project.Name, repository.Url)
				}
			}
		}
	}
	wg.Wait()
}
