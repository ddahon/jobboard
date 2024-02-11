package collectors

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/ddahon/jobboard/internal/models"
	"github.com/gocolly/colly/v2"
)

var companyShortname = "datadog"
var baseDomain = "careers.datadoghq.com"
var baseUrl = "https://" + baseDomain
var company *models.Company

func ScrapeDatadog() ([]models.Job, error) {
	company = models.GetCompanyByShortname(companyShortname)
	if company == nil {
		return nil, errors.New("Cannot retrieve company for shortname " + companyShortname + ". Aborting scraping for this company.")
	}
	log.Printf("company: %v", company)
	ctx, cancel := chromedp.NewContext(context.Background())

	defer cancel()

	jobLinks, err := getJobLinks(ctx)
	if err != nil {
		return nil, err
	}

	jobs, err := extractJobInfo(jobLinks)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func getJobLinks(ctx context.Context) ([]string, error) {
	jobLinks := make([]string, 0)

	if err := chromedp.Run(ctx, chromedp.Navigate(baseUrl+"/remote")); err != nil {
		return nil, err
	}

	for {
		jobLinks = append(jobLinks, getJobLinksFromCurrentPage(ctx)...)
		log.Println(len(jobLinks))
		err := chromedp.Run(ctx, chromedp.Click("a.ais-Pagination-link[aria-label=\"Next\"]", chromedp.AtLeast(0)))
		if err != nil {
			log.Println(err)
			break
		}
		time.Sleep(1 * time.Second)
	}

	return jobLinks, nil
}

func getJobLinksFromCurrentPage(ctx context.Context) []string {
	var nodes []*cdp.Node
	jobLinks := make([]string, 0)

	if err := chromedp.Run(ctx, chromedp.Nodes("button.job-card > a[href]", &nodes, chromedp.ByQueryAll)); err != nil {
		log.Println(err)
	} else {
		for _, node := range nodes {
			if href, found := node.Attribute("href"); found {
				jobLinks = append(jobLinks, baseUrl+href)
			}
		}
	}
	return jobLinks
}

func extractJobInfo(urls []string) ([]models.Job, error) {
	jobs := make([]models.Job, 0)
	c := colly.NewCollector(colly.AllowedDomains(baseDomain))

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.OnScraped(func(r *colly.Response) {
		jobs[len(jobs)-1].Link = r.Request.URL.String()
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to visit %v: %v", r.Request.URL, err)
	})
	c.OnHTML("main > h2", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Title = h.Text
	})
	c.OnHTML("main > h2 + div", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Location = h.Text
	})
	c.OnHTML(".job-description", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Description = h.Text
	})
	c.OnHTML(".job-description", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Description = h.Text
	})
	for _, url := range urls {
		job := models.Job{Company: company}
		jobs = append(jobs, job)
		c.Visit(url)
	}
	return jobs, nil
}
