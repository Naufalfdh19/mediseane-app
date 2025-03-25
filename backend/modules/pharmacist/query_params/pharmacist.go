package queryparams

import (
	"fmt"
	"strings"
)

type QueryParams struct {
	Name             string
	Email            string
	SipaNumber       string
	PhoneNumber      string
	YearOfExperience int
	SortBy           string
	Order            string
	Limit            int
	Page             int
}

type QueryParamsExistence struct {
	Name             string
	Email            string
	SipaNumber       string
	PhoneNumber      string
	YearOfExperience string
	SortBy           string
	Order            string
	Limit            string
	Page             string
	IsNameExists bool
	IsEmailExists bool
	IsSipaNumberExists bool
	IsPhoneNumberExists bool
	IsYearOfExperienceExists bool
	IsSortByExists bool
	IsOrderExists bool
	IsLimitExists bool
	IsPageExists bool
}


func AddQueryParams(params *[]any, queryParams QueryParams) string {
	var query string
	querIndex := 1

	if queryParams.Name != "" {
		query += fmt.Sprintf(" AND u.name ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.Name)
	}

	if queryParams.Email != "" {
		query += fmt.Sprintf(" AND u.email ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.Email)
	}

	if queryParams.SipaNumber != "" {
		query += fmt.Sprintf(" AND pd.sipa_number ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.SipaNumber)
	}

	if queryParams.PhoneNumber != "" {
		query += fmt.Sprintf(" AND pd.phone_number ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.PhoneNumber)
	}

	if queryParams.YearOfExperience != 0 {
		query += fmt.Sprintf(" AND pd.year_of_experience = $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.YearOfExperience)
	}

	query += sortByQuery(queryParams)

	if queryParams.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.Limit)
	}

	if queryParams.Page != 0 {
		query += fmt.Sprintf(" OFFSET %d", (queryParams.Page-1)*queryParams.Limit)
	}

	return query
}

func sortByQuery(queryParams QueryParams) string {
	query := ""
	if strings.ToLower(queryParams.Order) == "desc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by u.created_at desc"
		}
		if strings.ToLower(queryParams.SortBy) == "name" {
			query += " order by u.name desc"
		}
		if strings.ToLower(queryParams.SortBy) == "pharmacy" {
			query += " order by pd.pharmacy_id desc"
		}
	} else if strings.ToLower(queryParams.Order) == "asc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by u.created_at asc"
		}
		if strings.ToLower(queryParams.SortBy) == "name" {
			query += " order by u.name asc"
		}
		if strings.ToLower(queryParams.SortBy) == "pharmacy" {
			query += " order by pd.pharmacy_id asc "
		}
	}
	return query
}


func DefaultQuery(queryParams QueryParams) QueryParams {
	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}
	if queryParams.Page == 0 {
		queryParams.Page = 1
	}
	if queryParams.SortBy == "" {
		queryParams.SortBy = "created_at"
	}
	if queryParams.Order == "" {
		queryParams.Order = "desc"
	}
	return queryParams
}