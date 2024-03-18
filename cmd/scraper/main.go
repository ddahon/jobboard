package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ddahon/jobboard/cmd/scraper/analytics"

	"github.com/ddahon/jobboard/internal/pkg/models"
	"github.com/ddahon/jobboard/internal/scraper/companies/canonical"
	"github.com/ddahon/jobboard/internal/scraper/companies/datadog"
	"github.com/ddahon/jobboard/internal/scraper/companies/lumenalta"
	"github.com/ddahon/jobboard/internal/scraper/companies/spacelift"
	_ "github.com/lib/pq"
)

type ScrapeFunc func() ([]models.Job, error)

var allCollectors = map[string]ScrapeFunc{
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
		retries := 3
		jobs, fails := getJobs(retries, name, scrape)
		if jobs == nil {
			continue
		}
		nbFound := 0
		for _, job := range jobs {
			if err := job.Save(); err != nil {
				log.Printf("Failed to register job in DB: %v", err)
			}
			nbFound++
		}
		scrapeStats[name] = analytics.ScrapeResult{NbFound: nbFound, Failed: retries == fails, Retries: fails}
	}

	models.DeleteDeadJobs(scrapeStats)

	b, err := json.MarshalIndent(scrapeStats, "", "  ")
	if err == nil {
		log.Println(string(b))
	}
}

func getJobs(retries int, company string, scrape ScrapeFunc) ([]models.Job, int) {
	fails := 0
	var jobs []models.Job
	var err error
	log.Println("starting scraping of company:", company)
	for i := 0; i < retries; i++ {
		jobs, err = scrape()
		if err != nil {
			log.Printf("failed to scrape %v: %v\n", company, err)
			fails++
			continue
		}
		if jobs == nil {
			log.Printf("no jobs found for company %v", company)
			fails++
			continue
		}
	}
	return jobs, fails
}

func updateCollectorsList(args []string) {
	newCollectors := map[string]ScrapeFunc{}
	for _, e := range args {
		if val, ok := allCollectors[e]; ok {
			newCollectors[e] = val
		}
	}
	allCollectors = newCollectors
}
