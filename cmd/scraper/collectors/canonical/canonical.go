package canonical

import (
	"errors"

	"github.com/ddahon/jobboard/cmd/scraper/base"
	"github.com/ddahon/jobboard/internal/models"
)

var company *models.Company
var companyShortname = "canonical"

var CanonicalScraper = base.JsScraper[canonicalJob]{
	ScriptSelector: base.Selector{
		Type:  base.HTMLSelector,
		Value: `script[type="text/javascript"]`,
	},
	JsVarName: "vacancies",
	BaseScraper: base.BaseScraper{
		BaseDomain:  "canonical.com",
		StartingUrl: "https://canonical.com/careers/all",
	},
}

func Scrape() ([]models.Job, error) {
	company = models.GetCompanyByShortname(companyShortname)
	if company == nil {
		return nil, errors.New("Cannot retrieve company for shortname " + companyShortname + ". Aborting scraping for this company.")
	}
	return CanonicalScraper.Scrape()
}

type canonicalJob struct {
	Title       string
	Description string
	Location    string
	Link        string `json:"url"`
}

func (cj canonicalJob) ToJob() models.Job {
	return models.Job{
		Title:       cj.Title,
		Description: cj.Description,
		Link:        cj.Link,
		Location:    cj.Location,
		Company:     company,
	}
}
