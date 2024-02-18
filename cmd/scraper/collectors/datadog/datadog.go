package datadog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ddahon/jobboard/cmd/scraper/collectors"
	"github.com/ddahon/jobboard/internal/models"
)

var companyShortname = "datadog"
var baseDomain = "careers.datadoghq.com"
var baseUrl = "https://" + baseDomain
var startingUrl = baseUrl + "/remote"
var company *models.Company

func Scrape() ([]models.Job, error) {
	company = models.GetCompanyByShortname(companyShortname)
	if company == nil {
		return nil, errors.New("Cannot retrieve company for shortname " + companyShortname + ". Aborting scraping for this company.")
	}

	// get algolia tokens from js
	algoliaApiKey, algoliaAppId, algoliaIndex, err := getAlgoliaParams()
	if err != nil {
		return nil, err
	}
	return getJobs(algoliaApiKey, algoliaAppId, algoliaIndex)
}

func getAlgoliaParams() (string, string, string, error) {
	var apiKey, appId, index string
	jsSrc, err := collectors.FetchFileContent("https://careers.datadoghq.com/assets/scripts/main-YQ6Q5FDQ.js")
	if err != nil {
		return apiKey, appId, index, err
	}
	apiKey = collectors.GetJsKV(jsSrc, "ALGOLIA_PUBLIC_TOKEN")
	appId = collectors.GetJsKV(jsSrc, "ALGOLIA_APPLICATION")
	index = collectors.GetJsKV(jsSrc, "ALGOLIA_INDEX")

	return apiKey, appId, index, nil
}

type algoliaResult struct {
	Results []struct {
		Hits    []datadogJob
		Page    int
		NbPages int
	}
}

type datadogJob struct {
	Title       string
	Link        string `json:"absolute_url"`
	Location    string `json:"location_string"`
	Description string
}

func (dj *datadogJob) toJob() models.Job {
	return models.Job{
		Title:       dj.Title,
		Description: dj.Description,
		Location:    dj.Location,
		Link:        dj.Link,
		Company:     company,
	}
}

func getJobs(apiKey string, appId string, index string) ([]models.Job, error) {
	var jobs []models.Job
	algoliaUrl := fmt.Sprintf("https://%v-dsn.algolia.net/1/indexes/*/queries?&x-algolia-api-key=%v&x-algolia-application-id=%v", strings.ToLower(appId), apiKey, appId)
	page := 1
	for {
		reqBody := `{"requests":[{"indexName":"` + index + `","params":"facets=%5B%22time_type%22%2C%22parent_department_Engineering%22%2C%22child_department_Engineering%22%2C%22parent_department_Marketing%22%2C%22child_department_Marketing%22%2C%22parent_department_GeneralAdministrative%22%2C%22child_department_GeneralAdministrative%22%2C%22parent_department_TechnicalSolutions%22%2C%22child_department_TechnicalSolutions%22%2C%22parent_department_Sales%22%2C%22child_department_Sales%22%2C%22parent_department_ProductDesign%22%2C%22parent_department_ProductManagement%22%2C%22region_Americas%22%2C%22location_Americas%22%2C%22region_EMEA%22%2C%22location_EMEA%22%2C%22region_APAC%22%2C%22location_APAC%22%2C%22remote%22%5D&filters=remote%3A%20Remote%20AND%20(NOT%20ripplematch%3A%20true)&highlightPostTag=__%2Fais-highlight__&highlightPreTag=__ais-highlight__&hitsPerPage=10&maxValuesPerFacet=50&page=` + fmt.Sprint(page) + `&tagFilters="}]}`

		res, err := queryAlgolia(algoliaUrl, reqBody)
		if err != nil {
			log.Printf("failed to get page from algolia: %v", err)
			if page == 1 { // fail if first request fails
				return jobs, err
			}
		}
		if len(res.Results) < 1 {
			log.Printf("no results for page %v from algolia", page)
			if page == 1 { // fail if first request fails
				return jobs, nil
			}
		}
		for _, dj := range res.Results[0].Hits {
			jobs = append(jobs, dj.toJob())
		}

		page++
		if page > res.Results[0].NbPages {
			break
		}
	}
	return jobs, nil
}

func queryAlgolia(url string, body string) (algoliaResult, error) {
	result := algoliaResult{}
	res, err := http.Post(url, "application/json; charset=UTF-8", strings.NewReader(body))
	if err != nil {
		return result, err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}
	if !json.Valid([]byte(resBody)) {
		return result, fmt.Errorf("algolia query returned an invalid json response: %v", string(resBody))
	}
	if err := json.Unmarshal(resBody, &result); err != nil {
		return result, err
	}

	return result, nil
}
