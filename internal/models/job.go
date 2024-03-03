package models

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ddahon/jobboard/cmd/scraper/analytics"
)

type Job struct {
	Id          uint
	Description string
	Title       string
	Link        string
	Company     *Company
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Location    string // TODO make a dedicated Location struct
	Category    string // TODO make an enum ?
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
	return err != nil || res.StatusCode != 200
}

func DeleteDeadJobs(scrapeStats map[string]analytics.ScrapeResult) {
	log.Println("Cleaning up outdated jobs")
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
		if e, ok := scrapeStats[*j.Company.Shortname]; ok {
			e.NbDeleted++
			scrapeStats[*j.Company.Shortname] = e
		}
	}
}

func GetAllJobs() ([]Job, error) {
	var jobs []Job
	res, err := db.Query("SELECT * FROM jobs")
	if err != nil {
		return nil, err
	}
	defer res.Close()
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
	defer res.Close()
	if !res.Next() { // no results found
		return false, id
	}
	res.Scan(&id)
	return true, id
}

func NewJob(rows *sql.Rows) (Job, error) {
	var j Job
	var companyId uint
	var description, title, link, location, category sql.NullString
	if err := rows.Scan(&j.Id, &description, &title, &link, &companyId, &j.CreatedAt, &j.UpdatedAt, &location, &category); err != nil {
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
	if location.Valid {
		j.Location = location.String
	}
	if category.Valid {
		j.Category = category.String
	}
	j.Company = GetCompanyById(companyId)

	return j, nil
}
