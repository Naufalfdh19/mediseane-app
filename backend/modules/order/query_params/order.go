package queryparams

import (
	"fmt"
	"strings"
)

type QueryParams struct {
	Status string
	Limit  int
	Page   int
	SortBy string
	Order  string
}

type QueryParamsDto struct {
	Status string `form:"status"`
	SortBy string `form:"sort_by"`
	Order  string `form:"order"`
	Limit  int    `form:"limit"`
	Page   int    `form:"page"`
}

func AddQueryParams(params *[]any, queryParams QueryParams) string {
	var query string
	querIndex := 2

	if queryParams.Status != "" {
		queryParams.Status = "%" + queryParams.Status + "%"
		query += fmt.Sprintf(" AND od.status ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.Status)
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
	query := ``
	if strings.ToLower(queryParams.Order) == "desc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by o.created_at desc"
		}
		if strings.ToLower(queryParams.SortBy) == "status" {
			query += " order by od.status desc"
		}
	} else if strings.ToLower(queryParams.Order) == "asc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by o.created_at asc"
		}
		if strings.ToLower(queryParams.SortBy) == "status" {
			query += " order by od.status asc"
		}
	}

	return query
}
