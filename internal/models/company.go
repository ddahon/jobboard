package models

import "errors"

type Company struct {
	Id          uint
	Name        string
	Description string
	Website     string
}

func GetCompany(name string) (Company, error) {
	var company Company
	res, err := DB.Query("SELECT * FROM companies WHERE shortname='$1'", name)
	if err != nil {
		return company, err
	}

	if !res.Next() {
		return company, errors.New("Company not found: " + name)
	}
	res.Scan(&company.Id, &company.Name, &company.Description, &company.Website)
	return company, nil
}
