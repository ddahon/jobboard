package spacelift

import (
	"errors"

	"github.com/ddahon/jobboard/internal/pkg/models"
	"github.com/ddahon/jobboard/internal/scraper"
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
	bs := scraper.BaseScraper{BaseDomain: "careers.spacelift.io", StartingUrl: "https://careers.spacelift.io/jobs", Company: company, CompanyShortname: companyShortname}
	jobLinks, err := scraper.GetJobLinks(scraper.Selector{Type: scraper.HTMLSelector, Value: "#jobs_list_container a[href]"}, bs)
	if err != nil {
		return nil, err
	}
	locSelector := scraper.Selector{Type: scraper.HTMLSelector, Value: "dt:contains('Remote status') + dd"}
	descSelector := scraper.Selector{Type: scraper.HTMLSelector, Value: "main > section:nth-of-type(2)"}
	titleSelector := scraper.Selector{Type: scraper.HTMLSelector, Value: "h1.font-company-header"}
	return scraper.ExtractJobInfo(jobLinks, locSelector, descSelector, titleSelector, bs)
}
