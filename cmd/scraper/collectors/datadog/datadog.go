package datadog

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/ddahon/jobboard/internal/models"
)

var companyShortname = "datadog"
var baseDomain = "careers.datadoghq.com"
var baseUrl = "https://" + baseDomain
var startingUrl = baseUrl + "/remote"
var company *models.Company

func Scrape() ([]models.Job, error) {
	var jobs []models.Job
	company = models.GetCompanyByShortname(companyShortname)
	if company == nil {
		return nil, errors.New("Cannot retrieve company for shortname " + companyShortname + ". Aborting scraping for this company.")
	}

	// get algolia tokens from js
	algoliaApiKey, algoliaAppId, algoliaIndex, err := getAlgoliaParams()
	if err != nil {
		return nil, err
	}
	log.Printf("%v, %v, %v", algoliaApiKey, algoliaAppId, algoliaIndex)

	return jobs, nil
}

func getAlgoliaParams() (string, string, string, error) {
	var apiKey, appId, index string
	jsSource := "https://careers.datadoghq.com/assets/scripts/main-YQ6Q5FDQ.js"
	res, err := http.Get(jsSource)
	if err != nil {
		log.Printf("Failed to get js file: %v", err)
		return apiKey, appId, index, err
	}
	body, err := io.ReadAll(res.Body)
	sb := string(body)
	if err != nil {
		return apiKey, appId, index, err
	}
	if res.StatusCode != 200 {
		return apiKey, appId, index, fmt.Errorf("failed to get %v. response: %v", jsSource, sb)
	}
	r, err := regexp.Compile(`ALGOLIA_APPLICATION:"([^"]+)"`)
	if err != nil {
		return apiKey, appId, index, err
	}
	if m := r.FindStringSubmatch(sb); m != nil {
		appId = m[1]
	}
	r, err = regexp.Compile(`ALGOLIA_PUBLIC_TOKEN:"([^"]+)"`)
	if err != nil {
		return apiKey, appId, index, err
	}
	if m := r.FindStringSubmatch(sb); m != nil {
		apiKey = m[1]
	}
	r, err = regexp.Compile(`ALGOLIA_INDEX:"([^"]+)"`)
	if err != nil {
		return apiKey, appId, index, err
	}
	if m := r.FindStringSubmatch(sb); m != nil {
		index = m[1]
	}

	return apiKey, appId, index, nil
}
