package canonical

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/ddahon/jobboard/cmd/scraper/collectors"
	"github.com/ddahon/jobboard/internal/models"
	"github.com/gocolly/colly/v2"
)

var companyShortname = "canonical"
var baseDomain = "canonical.com"
var baseUrl = "https://" + baseDomain
var startingUrl = baseUrl + "/careers/all"
var company *models.Company

func Scrape() ([]models.Job, error) {
	var jobs []models.Job
	company = models.GetCompanyByShortname(companyShortname)
	if company == nil {
		return nil, errors.New("Cannot retrieve company for shortname " + companyShortname + ". Aborting scraping for this company.")
	}

	c := colly.NewCollector(colly.AllowedDomains(baseDomain))
	c.OnHTML(`script[type="text/javascript"]`, func(h *colly.HTMLElement) {
		jsSrc := h.Text
		v, err := collectors.GetJsArrayVar(jsSrc, "vacancies")
		if err != nil {
			log.Printf("Error while parsing js script: %v", err)
			return
		}
		var res []canonicalJob
		if err := json.Unmarshal([]byte(v), &res); err != nil {
			log.Printf("Error while parsing js: %v", err)
		}
		for _, j := range res {
			jobs = append(jobs, j.toJob())
		}
	})
	err := c.Visit(startingUrl)
	if err != nil {
		return jobs, err
	}

	return jobs, nil
}

type canonicalJob struct {
	Title       string
	Description string
	Location    string
	Link        string `json:"url"`
}

func (cj *canonicalJob) toJob() models.Job {
	return models.Job{
		Title:       cj.Title,
		Description: cj.Description,
		Link:        cj.Link,
		Location:    cj.Location,
		Company:     company,
	}
}
