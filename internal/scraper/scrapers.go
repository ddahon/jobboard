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

// TODO: there has to be a better way
func registerCallbackText(c *colly.Collector, s Selector, callback func(string)) {
	if s.Type == HTMLSelector {
		c.OnHTML(s.Value, func(h *colly.HTMLElement) { callback(h.Text) })
	}
	if s.Type == XMLSelector {
		c.OnXML(s.Value, func(x *colly.XMLElement) { callback(x.Text) })
	}
}

func registerCallbackLink(c *colly.Collector, s Selector, callback func(string)) {
	if s.Type == HTMLSelector {
		c.OnHTML(s.Value, func(h *colly.HTMLElement) { callback(h.Attr("href")) })
	}
	if s.Type == XMLSelector {
		c.OnXML(s.Value, func(x *colly.XMLElement) { callback(x.Attr("href")) })
	}
}

func ScrapeJsVar[T JobAdapter](s Selector, name string, bs BaseScraper) ([]models.Job, error) {
	var jobs []models.Job
	c := colly.NewCollector(colly.AllowedDomains(bs.BaseDomain))
	registerCallbackText(c, s, func(s string) {
		v, err := GetJsArrayVar(s, name)
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
	err := c.Visit(bs.StartingUrl)

	return jobs, err
}

func GetJobLinks(s Selector, bs BaseScraper) ([]string, error) {
	var jobLinks []string
	c := colly.NewCollector(colly.AllowedDomains(bs.BaseDomain))
	registerCallbackLink(c, s, func(url string) {
		jobLinks = append(jobLinks, url)
	})
	err := c.Visit(bs.StartingUrl)
	return jobLinks, err
}

func ExtractJobInfo(urls []string, locSelector Selector, descSelector Selector, titleSelector Selector, bs BaseScraper) ([]models.Job, error) {
	jobs := make([]models.Job, 0)
	c := colly.NewCollector(colly.AllowedDomains(bs.BaseDomain))
	c.OnScraped(func(r *colly.Response) {
		jobs[len(jobs)-1].Link = r.Request.URL.String()
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to visit %v: %v", r.Request.URL, err)
	})
	c.OnHTML(locSelector.Value, func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Location = h.Text
	})
	registerCallbackText(c, locSelector, func(s string) { jobs[len(jobs)-1].Location = s })
	registerCallbackText(c, descSelector, func(s string) { jobs[len(jobs)-1].Description = s })
	registerCallbackText(c, titleSelector, func(s string) { jobs[len(jobs)-1].Title = s })

	for _, url := range urls {
		job := models.Job{Company: bs.Company}
		jobs = append(jobs, job)
		if err := c.Visit(url); err != nil {
			log.Printf("Error while extracting job info in %v: %v", url, err)
		}
	}
	return jobs, nil
}
