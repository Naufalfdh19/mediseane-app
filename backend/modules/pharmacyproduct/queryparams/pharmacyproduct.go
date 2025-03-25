package queryparams

import (
	"fmt"
	"strings"
)

type QueryParamsDto struct {
	Name                  string `form:"name"`
	GenericName           string `form:"generic_name"`
	Manufacture           string `form:"manufacture"`
	ProductClassification string `form:"product_classification"`
	ProductForm           string `form:"product_form"`
	IsActive              string `form:"is_active"`
	SortBy                string `form:"sort_by"`
	Order                 string `form:"order"`
	Limit                 int    `form:"limit"`
	Page                  int    `form:"page"`
}

type QueryParams struct {
	Name                  string
	GenericName           string
	Manufacture           string
	ProductClassification string
	ProductForm           string
	IsActive              string
	SortBy                string
	Order                 string
	Limit                 int
	Page                  int
}

func AddQueryParams(params *[]any, queryParams QueryParams, querIndex *int) string {
	var query string

	if queryParams.Name != "" {
		queryParams.Name = "%" + queryParams.Name + "%"
		query += fmt.Sprintf(" AND p.name ILIKE $%d", *querIndex)
		*querIndex++
		*params = append(*params, queryParams.Name)
	}

	if queryParams.GenericName != "" {
		queryParams.GenericName = "%" + queryParams.GenericName + "%"
		query += fmt.Sprintf(" AND p.generic_name ILIKE $%d", *querIndex)
		*querIndex++
		*params = append(*params, queryParams.GenericName)
	}

	if queryParams.Manufacture != "" {
		queryParams.Manufacture = "%" + queryParams.Manufacture + "%"
		query += fmt.Sprintf(" AND p.manufacture ILIKE $%d", *querIndex)
		*querIndex++
		*params = append(*params, queryParams.Manufacture)
	}

	if queryParams.ProductClassification != "" {
		queryParams.ProductClassification = "%" + queryParams.ProductClassification + "%"
		query += fmt.Sprintf(" AND pc.name ILIKE $%d", *querIndex)
		*querIndex++
		*params = append(*params, queryParams.ProductClassification)
	}

	if queryParams.ProductForm != "" {
		queryParams.ProductForm = "%" + queryParams.ProductForm + "%"
		query += fmt.Sprintf(" AND pf.name ILIKE $%d", *querIndex)
		*querIndex++
		*params = append(*params, queryParams.ProductForm)
	}

	if queryParams.IsActive != "" && (queryParams.IsActive == "true" || queryParams.IsActive == "false") {
		query += fmt.Sprintf(" AND pp.is_active = $%d", *querIndex)
		*querIndex++
		*params = append(*params, queryParams.IsActive)
	}

	return query
}

func AddPaginationQuery(params *[]any, queryParams QueryParams, querIndex *int) string {
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

func AddSortByQuery(queryParams QueryParams) string {
	query := ""
	if strings.ToLower(queryParams.Order) == "desc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by pp.created_at desc"
		}
		if strings.ToLower(queryParams.SortBy) == "name" {
			query += " order by p.name desc"
		}
		if strings.ToLower(queryParams.SortBy) == "stock" {
			query += " order by pp.stock desc"
		}
	} else if strings.ToLower(queryParams.Order) == "asc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by pp.created_at asc"
		}
		if strings.ToLower(queryParams.SortBy) == "name" {
			query += " order by p.name asc"
		}
		if strings.ToLower(queryParams.SortBy) == "stock" {
			query += " order by pp.stock asc"
		}
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