package models

import (
	"database/sql"
	"log"
)

type Job struct {
	Id          uint
	Title       string
	Description string
	Link        string
	Company     *Company
	Languages   []string
	SalaryBegin uint
	SalaryEnd   uint
}

func AddJob(job Job) error {
	_, err := db.Exec("INSERT INTO jobs (description, title, link, company_id) VALUES ($1, $2, $3, $4)", job.Description, job.Title, job.Link, job.Company.Id)
	return err
}

func GetAllJobs() ([]Job, error) {
	var allJobs []Job
	for _, c := range allCompanies {
		jobs, err := c.GetAllJobs()
		if err != nil {
			log.Printf("Failed to retrieve jobs for company %v: %v", c.Shortname, err)
		}
		allJobs = append(allJobs, jobs...)
	}

	return allJobs, nil
}

func (j *Job) Scan(rows *sql.Rows) error {
	var companyId uint
	if err := rows.Scan(&j.Id, &j.Description, &j.Title, &j.Link, &companyId); err != nil {
		return err
	}
	j.Company = GetCompanyById(companyId)
	return nil
}
