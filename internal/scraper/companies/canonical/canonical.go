package canonical

import (
	"errors"

	"github.com/ddahon/jobboard/internal/pkg/models"
	"github.com/ddahon/jobboard/internal/scraper"
)

var company *models.Company
var companyShortname = "canonical"

func Scrape() ([]models.Job, error) {
	company = models.GetCompanyByShortname(companyShortname)
	if company == nil {
		return nil, errors.New("Cannot retrieve company for shortname " + companyShortname + ". Aborting scraping for this company.")
	}
	return scraper.ScrapeJsVar[canonicalJob](
		scraper.Selector{
			Type:  scraper.HTMLSelector,
			Value: `script[type="text/javascript"]`,
		},
		"vacancies",
		scraper.BaseScraper{
			BaseDomain:  "canonical.com",
			StartingUrl: "https://canonical.com/careers/all",
		},
	)
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
