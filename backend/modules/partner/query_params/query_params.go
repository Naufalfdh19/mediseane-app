package queryparams

import (
	"fmt"
	"strings"
)

type QueryParams struct {
	Name        string
	YearFounded string
	ActiveDays  string
	StartHour   string
	EndHour     string
	IsActive    string
	SortBy      string
	Order       string
	Limit       int
	Page        int
}

type QueryParamsExistence struct {
	Name                string
	YearFounded         string
	ActiveDays          string
	StartHour           string
	EndHour             string
	IsActive            string
	SortBy              string
	Order               string
	Limit               string
	Page                string
	IsNameExists        bool
	IsYearFoundedExists bool
	IsActiveDaysExists  bool
	IsStartHourExists   bool
	IsEndHourExists     bool
	IsActiveExists      bool
	IsSortByExsts       bool
	IsOrderExists       bool
	IsLimitExists       bool
	IsPageExists        bool
}

func AddQueryParams(params *[]any, queryParams QueryParams) string {
	var query string
	querIndex := 1

	if queryParams.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.Name)
	}

	if queryParams.YearFounded != "" {
		query += fmt.Sprintf(" AND year_founded ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.YearFounded)
	}

	if queryParams.ActiveDays != "" {
		query += fmt.Sprintf(" AND active_days ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.ActiveDays)
	}

	if queryParams.StartHour != "" {
		query += fmt.Sprintf(" AND start_hour ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.StartHour)
	}

	if queryParams.EndHour != "" {
		query += fmt.Sprintf(" AND end_hour ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.EndHour)
	}

	if queryParams.IsActive != "" {
		query += fmt.Sprintf(" AND is_active = $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.IsActive)
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
			query += " order by created_at desc"
		}
		if strings.ToLower(queryParams.SortBy) == "name" {
			query += " order by name desc"
		}
		if strings.ToLower(queryParams.SortBy) == "year_founded" {
			query += " order by year_founded desc"
		}

	} else if strings.ToLower(queryParams.Order) == "asc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by created_at asc"
		}
		if strings.ToLower(queryParams.SortBy) == "name" {
			query += " order by name asc"
		}
		if strings.ToLower(queryParams.SortBy) == "year_founded" {
			query += " order by year_founded asc"
		}
	}
	return query
}

