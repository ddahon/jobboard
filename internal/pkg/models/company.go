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
	defer res.Close()

	for res.Next() {
		job, err := NewJob(res)
		if err != nil {
			log.Printf("Failed to retrieve job: %v", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
