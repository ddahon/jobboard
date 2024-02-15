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

	return extractJobInfo(jobLinks)
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

func extractJobInfo(urls []string) ([]models.Job, error) {
	jobs := make([]models.Job, 0)
	c := colly.NewCollector(colly.AllowedDomains(baseDomain))
	c.OnScraped(func(r *colly.Response) {
		jobs[len(jobs)-1].Link = r.Request.URL.String()
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to visit %v: %v", r.Request.URL, err)
	})
	c.OnHTML("dt:contains('Remote status') + dd", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Location = h.Text
	})
	c.OnHTML("main > section:nth-of-type(2)", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Description = h.Text
	})
	c.OnHTML("h1.font-company-header", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Title = h.Text
	})

	for _, url := range urls {
		job := models.Job{Company: company}
		jobs = append(jobs, job)
		if err := c.Visit(url); err != nil {
			log.Printf("Error while extracting job info in %v: %v", url, err)
		}
	}
	return jobs, nil
}
