package base

import (
	"encoding/json"
	"log"

	"github.com/ddahon/jobboard/internal/models"
	"github.com/gocolly/colly/v2"
)

type JobAdapter interface {
	ToJob() models.Job
}

type JsScraper[T JobAdapter] struct {
	ScriptSelector Selector // to get the script src
	JsVarName      string   // name of JS var that contains the jobs
	BaseScraper
}

func (s *JsScraper[T]) Scrape() ([]models.Job, error) {
	var jobs []models.Job

	c := colly.NewCollector(colly.AllowedDomains(s.BaseDomain))
	if s.ScriptSelector.Type == HTMLSelector {
		c.OnHTML(s.ScriptSelector.Value, func(h *colly.HTMLElement) {
			jsSrc := h.Text
			v, err := GetJsArrayVar(jsSrc, s.JsVarName)
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
	err := c.Visit(s.StartingUrl)

	return jobs, err
}
