package lumenalta

import (
	"errors"
	"log"

	"github.com/ddahon/jobboard/internal/pkg/models"
	"github.com/gocolly/colly/v2"
)

var companyShortname = "lumenalta"
var baseDomain = "lumenalta.com"
var baseUrl = "https://" + baseDomain
var company *models.Company

func Scrape() ([]models.Job, error) {
	company = models.GetCompanyByShortname(companyShortname)
	if company == nil {
		return nil, errors.New("Cannot retrieve company for shortname " + companyShortname + ". Aborting scraping for this company.")
	}
	jobLinks, err := getJobLinks()
	if err != nil {
		return nil, err
	}

	return extractJobInfo(jobLinks)
}

func getJobLinks() ([]string, error) {
	var jobLinks []string
	c := colly.NewCollector(colly.AllowedDomains(baseDomain))
	c.OnXML("//a[contains(@href, \"/jobs/\")]", func(x *colly.XMLElement) {
		jobLinks = append(jobLinks, baseUrl+"/remote-jobs/"+x.Attr("href"))
	})
	err := c.Visit(baseUrl + "/remote-jobs")
	return jobLinks, err
}

func extractJobInfo(urls []string) ([]models.Job, error) {
	jobs := make([]models.Job, 0)
	c := colly.NewCollector(colly.AllowedDomains(baseDomain))
	c.OnScraped(func(r *colly.Response) {
		jobs[len(jobs)-1].Link = r.Request.URL.String()
		jobs[len(jobs)-1].Location = "Remote"
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to visit %v: %v", r.Request.URL, err)
	})
	c.OnXML("//h4[contains(text(), 'About the Role')]/following::p[1]", func(x *colly.XMLElement) {
		jobs[len(jobs)-1].Description = x.Text
	})
	c.OnXML("//h1[contains(@class, 'hero-header')]", func(x *colly.XMLElement) {
		jobs[len(jobs)-1].Title = x.Text
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
