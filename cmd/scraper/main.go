package main

import (
	"log"
	"os"

	"github.com/ddahon/jobboard/cmd/scraper/collectors/datadog"
	"github.com/ddahon/jobboard/cmd/scraper/collectors/spacelift"
	"github.com/ddahon/jobboard/internal/models"
	_ "github.com/lib/pq"
)

var allCollectors = map[string]func() ([]models.Job, error){
	"datadog":   datadog.Scrape,
	"spacelift": spacelift.Scrape,
}

func main() {
	if len(os.Args) > 1 {
		updateCollectorsList(os.Args)
	}

	connStr := "postgresql://postgres:password@localhost:5432/jobs?sslmode=disable"
	if err := models.InitDB(connStr); err != nil {
		log.Fatalln(err)
	}

	for name, scrape := range allCollectors {
		log.Println("starting scraping of company:", name)
		jobs, err := scrape()
		if err != nil {
			log.Printf("failed to scrape %v: %v\n", name, err)
			continue
		}
		if jobs == nil {
			log.Printf("no jobs found for company %v", name)
			continue
		}
		for _, job := range jobs {
			if err := job.Save(); err != nil {
				log.Printf("Failed to register job in DB: %v", err)
			}
		}
	}

	models.DeleteDeadJobs()
}

func updateCollectorsList(args []string) {
	newCollectors := map[string]func() ([]models.Job, error){}
	for _, e := range args {
		if val, ok := allCollectors[e]; ok {
			newCollectors[e] = val
		}
	}
	allCollectors = newCollectors
}
