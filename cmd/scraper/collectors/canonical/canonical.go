package canonical

import (
	"context"
	"errors"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/ddahon/jobboard/internal/models"
	"github.com/gocolly/colly/v2"
)

var companyShortname = "canonical"
var baseDomain = "canonical.com"
var baseUrl = "https://" + baseDomain
var company *models.Company

func Scrape() ([]models.Job, error) {
	company = models.GetCompanyByShortname(companyShortname)
	if company == nil {
		return nil, errors.New("Cannot retrieve company for shortname " + companyShortname + ". Aborting scraping for this company.")
	}
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
	var nodes []*cdp.Node
	jobLinks := make([]string, 0)
	err := chromedp.Run(ctx, chromedp.Navigate(baseUrl+"/careers/all?filter=Engineering"))
	if err != nil {
		return nil, err
	}
	log.Printf("hihi")
	if err := chromedp.Run(ctx, chromedp.Nodes("h3 a[href]", &nodes, chromedp.ByQueryAll)); err != nil {
		log.Println(err)
	} else {
		for _, node := range nodes {
			if href, found := node.Attribute("href"); found {
				jobLinks = append(jobLinks, baseUrl+href)
			}
		}
	}

	return jobLinks, nil
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
	c.OnHTML("h1", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Title = h.Text
	})
	c.OnHTML("p.p-muted-heading", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Location = h.Text
	})
	c.OnHTML("#details > div:nth-child(2)", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Description = h.Text
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
