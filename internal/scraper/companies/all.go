package companies

import (
	"github.com/ddahon/jobboard/internal/pkg/models"
	"github.com/ddahon/jobboard/internal/scraper/companies/canonical"
	"github.com/ddahon/jobboard/internal/scraper/companies/datadog"
	"github.com/ddahon/jobboard/internal/scraper/companies/lumenalta"
	"github.com/ddahon/jobboard/internal/scraper/companies/spacelift"
)

type ScrapeFunc func() ([]models.Job, error)

var AllCollectors = map[string]ScrapeFunc{
	"datadog":   datadog.Scrape,
	"spacelift": spacelift.Scrape,
	"canonical": canonical.Scrape,
	"lumenalta": lumenalta.Scrape,
}
