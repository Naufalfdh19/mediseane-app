package queryparams

import (
	"fmt"
	"strings"
)

type QueryParams struct {
	Name       string
	Email      string
	Role       string
	SortBy     string
	Order      string
	Limit      int
	Page       int
}

type QueryParamsExistence struct {
	Name       string
	Email      string
	Role       string
	SortBy     string
	Order      string
	Limit      string
	Page       string
	IsNameExists       bool
	IsEmailExists      bool
	IsRoleExists       bool
	IsSortByExsts      bool
	IsOrderExists      bool
	IsLimitExists      bool
	IsPageExists       bool
}

func AddQueryParams(params *[]any, queryParams QueryParams) string {
	var query string
	querIndex := 1

	if queryParams.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.Name)
	}

	if queryParams.Email != "" {
		query += fmt.Sprintf(" AND email ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.Email)
	}
	if queryParams.Role != "" {
		query += fmt.Sprintf(" AND role ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.Role)
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
	} else if strings.ToLower(queryParams.Order) == "asc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by u.created_at asc"
		}
		if strings.ToLower(queryParams.SortBy) == "name" {
			query += " order by u.name asc"
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