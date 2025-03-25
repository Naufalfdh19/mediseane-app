package queryparams

import (
	"fmt"
	"strings"
)


func AddMasterConditionQuery(params *[]any, queryParams QueryParams, querIndex *int) string {
	var query string

	if queryParams.Name != "" {
		queryParams.Name = "%" + queryParams.Name + "%"
		query += fmt.Sprintf(" AND p.name ILIKE $%d", *querIndex)
		*querIndex++
		*params = append(*params, queryParams.Name)
	}
	return query
}

func AddMasterSortByQuery(queryParams QueryParams) string {
	query := ``
	if strings.ToLower(queryParams.Order) == "desc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by p.created_at desc"
		}
		if strings.ToLower(queryParams.SortBy) == "name" {
			query += " order by p.name desc"
		}
	} else if strings.ToLower(queryParams.Order) == "asc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by p.created_at asc"
		}
		if strings.ToLower(queryParams.SortBy) == "name" {
			query += " order by p.name asc"
		}
	}

	return query
}

func AddMasterPaginationQuery(params *[]any, queryParams QueryParams, querIndex *int) string {
	query := ``
	if queryParams.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%d", *querIndex)
		*querIndex++
		*params = append(*params, queryParams.Limit)
	}

	if queryParams.Page != 0 {
		query += fmt.Sprintf(" OFFSET %d", (queryParams.Page-1)*queryParams.Limit)
	}

	return query
}

func DefaultQueryParams(queryParams *QueryParams) {
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
}
