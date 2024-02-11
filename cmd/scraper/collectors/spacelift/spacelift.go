package spacelift

import (
	"errors"
	"log"

	"github.com/ddahon/jobboard/internal/models"
	"github.com/gocolly/colly/v2"
)

var companyShortname = "spacelift"
var baseDomain = "careers.spacelift.io"
var baseUrl = "https://" + baseDomain
var company *models.Company

func Scrape() ([]models.Job, error) {
	company = models.GetCompanyByShortname(companyShortname)
	if company == nil {
		return nil, errors.New("Cannot retrieve company for shortname " + companyShortname + ". Aborting scraping for this company.")
	}
	log.Printf("company: %v", company)
	jobLinks, err := getJobLinks()
	if err != nil {
		return nil, err
	}
	log.Printf("jobLinks: %v", jobLinks)

	return nil, err
}

func getJobLinks() ([]string, error) {
	var jobLinks []string
	c := colly.NewCollector(colly.AllowedDomains(baseDomain))
	c.OnHTML("#jobs_list_container a[href]", func(h *colly.HTMLElement) {
		jobLinks = append(jobLinks, h.Attr("href"))
	})
	err := c.Visit(baseUrl + "/jobs")
	return jobLinks, err
}
