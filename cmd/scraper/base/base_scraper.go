package base

import "github.com/ddahon/jobboard/internal/models"

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
