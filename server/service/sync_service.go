package service

import (
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron"
	// "github.com/wadahiro/gitss/server/model"
	"github.com/wadahiro/gitss/server/repo"
)

var mutex = new(sync.Mutex)
var scheduler *cron.Cron

func RunSyncScheduler(reader *repo.GitRepoReader) {
	mutex.Lock()
	defer mutex.Unlock()

	if scheduler != nil {
		log.Println("Stop sync schduler...")
		scheduler.Stop()
	}
	scheduler = cron.New()


	scheduler.Start()
	fmt.Println("Started sync schduler...")
}