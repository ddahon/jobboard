package main

import (
	"log"

	"github.com/ddahon/jobboard/cmd/scraper/collectors"
	"github.com/ddahon/jobboard/internal/models"
	_ "github.com/lib/pq"
)

var allCollectors = map[string]func() ([]models.Job, error){
	"datadog": collectors.ScrapeDatadog,
}

func main() {
	connStr := "postgresql://postgres:password@localhost:5432/jobs?sslmode=disable"
	if err := models.InitDB(connStr); err != nil {
		log.Fatalln(err)
	}

	for name, scrape := range allCollectors {
		log.Println("starting scraping of company:", name)
		jobs, err := scrape()
		if err != nil || jobs == nil {
			log.Printf("failed to scrape %v: %v\n", name, err)
			continue
		}
		for _, job := range jobs {
			if err := job.Save(); err != nil {
				log.Printf("Failed to register job in DB: %v", err)
			}
		}
	}
}
