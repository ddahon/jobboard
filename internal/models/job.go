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
	Company     Company
	Languages   []string
	SalaryBegin uint
	SalaryEnd   uint
}

var DB *sql.DB

func AddJob(job Job) error {
	_, err := DB.Exec("INSERT INTO jobs (description, title, link, company_id) VALUES ($1, $2, $3)", job.Description, job.Title, job.Link, job.Company.Id)
	return err
}

func GetJobs() ([]Job, error) {
	res, err := DB.Query("SELECT * FROM jobs")
	if err != nil {
		return nil, err
	}

	var jobs []Job
	for res.Next() {
		var job Job
		var companyId uint
		if err := res.Scan(&job.Id, &job.Description, &job.Title, &job.Link, companyId); err != nil {
			log.Printf("Failed to retrieve row from DB: %v", err)
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
