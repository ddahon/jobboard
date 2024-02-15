package models

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"
)

type Job struct {
	Id          uint
	Title       string
	Description string
	Link        string
	Company     *Company
	Location    string // TODO make a dedicated Location struct
	Languages   []string
	SalaryBegin uint
	SalaryEnd   uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (j Job) Save() error {
	exists, id := CheckJobExists(j.Link)
	if exists {
		log.Printf("Job with link %v already exists. Updating existing definition", j.Link)
		_, err := db.Exec("UPDATE jobs SET description=$1, title=$2, location=$3, updated_at=$4 WHERE id=$5", strings.TrimSpace(j.Description), strings.TrimSpace(j.Title), strings.TrimSpace(j.Location), time.Now(), id)
		return err
	}
	_, err := db.Exec("INSERT INTO jobs (description, title, link, company_id, location) VALUES ($1, $2, $3, $4, $5)", strings.TrimSpace(j.Description), strings.TrimSpace(j.Title), strings.TrimSpace(j.Link), j.Company.Id, strings.TrimSpace(j.Location))
	return err
}

func (j Job) Delete() error {
	_, err := db.Exec("DELETE FROM jobs WHERE id=$1", j.Id)
	if err != nil {
		return err
	}
	log.Printf("Deleted job with link %v", j.Link)
	return nil
}

func isDeadLink(link string) bool {
	res, err := http.Get(link)
	return res.StatusCode != 200 || err != nil
}

func DeleteDeadJobs() {
	jobs, err := GetAllJobs()
	if err != nil {
		log.Printf("Could not retrieve jobs: %v. Skipping dead links checking", err)
	}
	for _, j := range jobs {
		if !isDeadLink(j.Link) {
			continue
		}
		if err := j.Delete(); err != nil {
			log.Printf("Could not delete job: %v", err)
		}
	}
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
			log.Printf("Failed to retrieve a job: %v", err)
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// CheckJobExists returns a boolean representing if the job exists and a  job id if it exists.
func CheckJobExists(link string) (bool, uint) {
	var id uint
	res, err := db.Query("SELECT id FROM jobs WHERE link=$1", link)
	if err != nil {
		log.Printf("Error while checking if job exists: %v", err)
		return true, id
	}
	if !res.Next() { // no results found
		return false, id
	}
	res.Scan(&id)
	return true, id
}

func NewJob(rows *sql.Rows) (Job, error) {
	var j Job
	var companyId uint
	var description, title, link sql.NullString
	if err := rows.Scan(&j.Id, &description, &title, &link, &companyId, &j.CreatedAt, &j.UpdatedAt, &j.Location); err != nil {
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
