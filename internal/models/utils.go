package models

import (
	"database/sql"
	"log"
)

var db *sql.DB
var allCompanies []Company

func InitDB(connStr string) error {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err

	}
	if err := initCompanies(); err != nil {
		return err
	}

	return nil
}

func initCompanies() error {
	res, err := db.Query("SELECT * FROM companies")
	if err != nil {
		return err
	}

	for res.Next() {
		var company Company
		err := res.Scan(&company.Id, &company.Name, &company.Description, &company.Website, &company.Shortname)
		if err != nil {
			log.Printf("Failed to retrieve company from DB: %v", err)
			continue
		}
		allCompanies = append(allCompanies, company)
		log.Printf("Got all companies: %v", allCompanies)
	}
	return nil
}
