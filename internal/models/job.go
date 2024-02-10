package models

import (
	"database/sql"
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
	var jobs []Job
	res, err := db.Query("SELECT * FROM jobs")
	if err != nil {
		return nil, err
	}
	for res.Next() {
		job, err := NewJob(res)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func NewJob(rows *sql.Rows) (Job, error) {
	var j Job
	var companyId uint
	var description, title, link sql.NullString
	if err := rows.Scan(&j.Id, &description, &title, &link, &companyId); err != nil {
		return j, err
	}
	if description.Valid {
		j.Description = description.String
	}
	if title.Valid {
		j.Title = title.String
	}
	if link.Valid {
		j.Link = link.String
	}
	j.Company = GetCompanyById(companyId)

	return j, nil
}
