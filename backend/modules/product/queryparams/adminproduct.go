package queryparams

import (
	"fmt"
	"strings"
)

type AdminQueryParamsDto struct {
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

type AdminQueryParams struct {
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

func AddAdminQueryParams(params *[]any, queryParams AdminQueryParams) string {
	var query string
	querIndex := 1

	if queryParams.Name != "" {
		queryParams.Name = "%" + queryParams.Name + "%"
		query += fmt.Sprintf(" AND p.name ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.Name)
	}

	if queryParams.GenericName != "" {
		queryParams.GenericName = "%" + queryParams.GenericName + "%"
		query += fmt.Sprintf(" AND p.generic_name ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.GenericName)
	}

	if queryParams.Manufacture != "" {
		queryParams.Manufacture = "%" + queryParams.Manufacture + "%"
		query += fmt.Sprintf(" AND p.manufacture ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.Manufacture)
	}

	if queryParams.ProductClassification != "" {
		queryParams.ProductClassification = "%" + queryParams.ProductClassification + "%"
		query += fmt.Sprintf(" AND pc.name ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.ProductClassification)
	}

	if queryParams.ProductForm != "" {
		queryParams.ProductForm = "%" + queryParams.ProductForm + "%"
		query += fmt.Sprintf(" AND pf.name ILIKE $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.ProductForm)
	}

	if queryParams.IsActive != "" && (queryParams.IsActive == "true" || queryParams.IsActive == "false") {
		query += fmt.Sprintf(" AND is_active = $%d", querIndex)
		querIndex++
		*params = append(*params, queryParams.IsActive)
	}

	query += AdminSortBy(queryParams)

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

func AdminSortBy(queryParams AdminQueryParams) string {
	query := ""
	if strings.ToLower(queryParams.Order) == "desc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by p.created_at desc"
		}
		if strings.ToLower(queryParams.SortBy) == "name" {
			query += " order by p.name desc"
		}
		if strings.ToLower(queryParams.SortBy) == "product_used" {
			query += " order by product_used desc"
		}
	} else if strings.ToLower(queryParams.Order) == "asc" {
		if strings.ToLower(queryParams.SortBy) == "created_at" {
			query += " order by p.created_at asc"
		}
		if strings.ToLower(queryParams.SortBy) == "name" {
			query += " order by p.name asc"
		}
		if strings.ToLower(queryParams.SortBy) == "product_used" {
			query += " order by product_used asc"
		}
	}
	return query
}

func AdminGetCategoriesParams(params *[]any, categories []int, queryIndex *int) string {
	query := ``
	for _, category := range categories {
		query += fmt.Sprintf(` OR id = $%d`, *queryIndex)
		*params = append(*params, category)
		*queryIndex++
	}
	query = strings.TrimPrefix(query, " OR")
	return "("+query+")"
}


func AdminDefaultQuery(queryParams AdminQueryParams) AdminQueryParams {
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


