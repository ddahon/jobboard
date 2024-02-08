package models

import "database/sql"

type Company struct {
	Id          uint
	Name        string
	Description string
	Website     string
}

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
	_, err := DB.Exec("INSERT INTO jobs (description, title, link) VALUES ($1, $2, $3)", job.Description, job.Title, job.Link)
	return err
}
