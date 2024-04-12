package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/ddahon/jobboard/cmd/scraper/analytics"

	"github.com/ddahon/jobboard/internal/pkg/models"
	"github.com/ddahon/jobboard/internal/scraper/companies"
)

func main() {
	var clean bool
	var sqliteDB string
	var collectors map[string]companies.ScrapeFunc

	collectors, clean, sqliteDB = parseArguments()

	if err := models.InitDB(sqliteDB); err != nil {
		log.Fatalln(err)
	}

	scrapeStats := map[string]analytics.ScrapeResult{}

	for name, scrape := range collectors {
		retries := 3
		jobs, fails := executeScrape(retries, name, scrape)
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

func parseArguments() (map[string]companies.ScrapeFunc, bool, string) {
	var clean bool
	var sqliteDB string
	flag.BoolVar(&clean, "clean", false, "delete dead jobs")
	flag.StringVar(&sqliteDB, "sqlite-db", "", "SQLite database connection string")
	flag.Parse()

	collectors := companies.AllCollectors
	if len(flag.Args()) > 0 {
		collectors = updateCollectorsList(flag.Args())
	}

	if sqliteDB == "" {
		log.Fatal("Please provide the SQLite database connection string using the -sqlite-db flag")
	}

	return collectors, clean, sqliteDB
}

func executeScrape(retries int, company string, scrape companies.ScrapeFunc) ([]models.Job, int) {
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
		break
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
