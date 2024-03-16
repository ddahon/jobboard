package lumenalta

import (
	"errors"

	"github.com/ddahon/jobboard/internal/pkg/models"
	"github.com/ddahon/jobboard/internal/scraper"
)

var companyShortname = "lumenalta"
var company *models.Company

func Scrape() ([]models.Job, error) {
	company = models.GetCompanyByShortname(companyShortname)
	if company == nil {
		return nil, errors.New("Cannot retrieve company for shortname " + companyShortname + ". Aborting scraping for this company.")
	}

	bs := scraper.BaseScraper{BaseDomain: "lumenalta.com", StartingUrl: "https://lumenalta.com/remote-jobs", Company: company, CompanyShortname: companyShortname}
	jobLinks, err := scraper.GetJobLinks(scraper.Selector{Type: scraper.XMLSelector, Value: "//a[contains(@href, \"/jobs/\")]"}, bs)
	if err != nil {
		return nil, err
	}
	locSelector := scraper.Selector{}
	descSelector := scraper.Selector{Type: scraper.HTMLSelector, Value: "//h4[contains(text(), 'About the Role')]/following::p[1]"}
	titleSelector := scraper.Selector{Type: scraper.HTMLSelector, Value: "//h1[contains(@class, 'hero-header')]"}

	return scraper.ExtractJobInfo(jobLinks, locSelector, descSelector, titleSelector, bs)
}
