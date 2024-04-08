package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/ddahon/jobboard/cmd/scraper/analytics"

	"github.com/ddahon/jobboard/internal/pkg/models"
	"github.com/ddahon/jobboard/internal/scraper/companies"
)

func main() {
	var clean bool
	flag.BoolVar(&clean, "clean", false, "delete dead jobs")
	flag.Parse()

	collectors := companies.AllCollectors
	if len(flag.Args()) > 0 {
		collectors = updateCollectorsList(flag.Args())
	}

	connStr := os.Getenv("SQLITE_DB")
	if err := models.InitDB(connStr); err != nil {
		log.Fatalln(err)
	}

	scrapeStats := map[string]analytics.ScrapeResult{}

	for name, scrape := range collectors {
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

	if clean {
		models.DeleteDeadJobs(scrapeStats)
	}

	b, err := json.MarshalIndent(scrapeStats, "", "  ")
	if err == nil {
		log.Println(string(b))
	}
}

func getJobs(retries int, company string, scrape companies.ScrapeFunc) ([]models.Job, int) {
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

func updateCollectorsList(args []string) map[string]companies.ScrapeFunc {
	newCollectors := map[string]companies.ScrapeFunc{}
	for _, e := range args {
		if val, ok := companies.AllCollectors[e]; ok {
			newCollectors[e] = val
		}
	}
	return newCollectors
}
