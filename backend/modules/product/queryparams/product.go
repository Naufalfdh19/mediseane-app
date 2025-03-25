package queryparams

import (
	"fmt"
)

type QueryParams struct {
	Name       string
	CategoryID int
	Limit      int
	Page       int
	SortBy     string
	Order      string
}

type QueryParamsDto struct {
	Name       string `form:"name"`
	CategoryID int    `form:"category_id"`
	SortBy     string `form:"sort_by"`
	Order      string `form:"order"`
	Limit      int    `form:"limit"`
	Page       int    `form:"page"`
}

func AddConditionQuery(params *[]any, queryParams QueryParams, querIndex *int) string {
	var query string

	if queryParams.Name != "" {
		queryParams.Name = "%" + queryParams.Name + "%"
		query += fmt.Sprintf(" AND p.name ILIKE $%d", *querIndex)
		*querIndex++
		*params = append(*params, queryParams.Name)
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

func AddFilterByCategoryID(params *[]any, queryParams QueryParams, querIndex *int) string {
	var query string
	if queryParams.CategoryID > 0 {
		query = fmt.Sprintf(` JOIN product_multi_categories pmc 
		ON pmc.product_id = pp.product_id 
		AND pmc.product_category_id = $%d`,*querIndex)
		*querIndex++
		*params = append(*params, queryParams.CategoryID)
	}

	return query
} 
