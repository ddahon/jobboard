package scraper

import (
	"encoding/json"
	"log"

	"github.com/ddahon/jobboard/internal/pkg/models"
	"github.com/gocolly/colly/v2"
)

type JobAdapter interface {
	ToJob() models.Job
}

type BaseScraper struct {
	BaseDomain       string
	StartingUrl      string
	Company          *models.Company
	CompanyShortname string
}

type Selector struct {
	Type  SelectorType
	Value string
}

type SelectorType int

const (
	HTMLSelector SelectorType = iota
	XMLSelector
)

func ScrapeJsVar[T JobAdapter](s Selector, name string, bs BaseScraper) ([]models.Job, error) {
	var jobs []models.Job

	c := colly.NewCollector(colly.AllowedDomains(bs.BaseDomain))
	if s.Type == HTMLSelector {
		c.OnHTML(s.Value, func(h *colly.HTMLElement) {
			jsSrc := h.Text
			v, err := GetJsArrayVar(jsSrc, name)
			if err != nil {
				log.Printf("Error while parsing js script: %v", err)
				return
			}
			var res []T
			if err := json.Unmarshal([]byte(v), &res); err != nil {
				log.Printf("Error while parsing js: %v", err)
			}
			for _, j := range res {
				jobs = append(jobs, j.ToJob())
			}
		})
	}
	err := c.Visit(bs.StartingUrl)

	return jobs, err
}