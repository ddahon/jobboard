package models

import "log"

type Company struct {
	Id          uint
	Name        *string
	Description *string
	Website     *string
	Shortname   *string
}

func GetCompanyByShortname(shortname string) *Company {
	for _, c := range allCompanies {
		if c.Shortname == nil {
			continue
		}
		if *c.Shortname == shortname {
			return &c
		}
	}
	return nil
}

func GetCompanyById(id uint) *Company {
	for _, c := range allCompanies {
		if c.Id == id {
			return &c
		}
	}
	return nil
}

func (c Company) GetAllJobs() ([]Job, error) {
	var jobs []Job
	res, err := db.Query("SELECT * FROM jobs WHERE company_id=$1", c.Id)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var job Job
		if err := job.Scan(res); err != nil {
			log.Printf("Failed to retrieve job: %v", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
