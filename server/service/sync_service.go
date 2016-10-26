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

func RunSyncScheduler(config config.Config, importer *importer.GitImporter) {
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

func RunSync(config config.Config, importer *importer.GitImporter) {
	list, err := config.GetAllIndexConf()
	if err != nil {
		log.Printf("Read conf error. %+v\n", err)
	}

	workers := util.GenWorkers(2)
	var wg sync.WaitGroup

	for i := range list {
		indexConf := list[i]
		log.Printf("Sync for %s:%s/%s\n", indexConf.Organization, indexConf.Repository, indexConf.Repository)

		wg.Add(1)
		workers <- func() {
			defer wg.Done()
			importer.Run(indexConf.Organization, indexConf.Project, indexConf.Url)
		}
	}
	wg.Wait()
}
