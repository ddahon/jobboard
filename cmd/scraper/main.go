package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ddahon/jobboard/cmd/scraper/analytics"
	"github.com/ddahon/jobboard/cmd/scraper/collectors/canonical"
	"github.com/ddahon/jobboard/cmd/scraper/collectors/datadog"
	"github.com/ddahon/jobboard/cmd/scraper/collectors/lumenalta"
	"github.com/ddahon/jobboard/cmd/scraper/collectors/spacelift"
	"github.com/ddahon/jobboard/internal/pkg/models"
	_ "github.com/lib/pq"
)

var allCollectors = map[string]func() ([]models.Job, error){
	"datadog":   datadog.Scrape,
	"spacelift": spacelift.Scrape,
	"canonical": canonical.Scrape,
	"lumenalta": lumenalta.Scrape,
}

func main() {
	if len(os.Args) > 1 {
		updateCollectorsList(os.Args)
	}

	connStr := "postgresql://postgres:password@localhost:5432/jobs?sslmode=disable"
	if err := models.InitDB(connStr); err != nil {
		log.Fatalln(err)
	}

	scrapeStats := map[string]analytics.ScrapeResult{}

	for name, scrape := range allCollectors {
		nbFound := 0
		log.Println("starting scraping of company:", name)
		jobs, err := scrape()
		if err != nil {
			log.Printf("failed to scrape %v: %v\n", name, err)
			scrapeStats[name] = analytics.ScrapeResult{Failed: true}
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
			nbFound++
		}
		scrapeStats[name] = analytics.ScrapeResult{NbFound: nbFound}
	}

	models.DeleteDeadJobs(scrapeStats)

	b, err := json.MarshalIndent(scrapeStats, "", "  ")
	if err == nil {
		log.Println(string(b))
	}
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
