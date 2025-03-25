package queryparams

import (
	"fmt"
	"montelukast/pkg/dateconverter"
)

type QueryParams struct {
	Name   string
	Limit  int
	Page   int
	SortBy string
	Order  string
}

type QueryParamsDto struct {
	Name   string `form:"name"`
	SortBy string `form:"sort_by"`
	Order  string `form:"order"`
	Limit  int    `form:"limit"`
	Page   int    `form:"page"`
}

func AddConditionQuery(params *[]any, queryParams QueryParams, querIndex *int) string {
	var query string

	if queryParams.Name != "" {
		queryParams.Name = "%" + queryParams.Name + "%"
		query += fmt.Sprintf(" AND product_name ILIKE $%d", *querIndex)
		*querIndex++
		*params = append(*params, queryParams.Name)
	}

	query += fmt.Sprintf(" AND start_hour < $%d", *querIndex)
	query += fmt.Sprintf(" AND end_hour > $%d", *querIndex)
	*params = append(*params, dateconverter.GetCurrentTime())	
	*querIndex++

	query += fmt.Sprintf(" AND active_days ilike $%d", *querIndex)
	*params = append(*params, "%" + dateconverter.GetCurrentDay() + "%")
	*querIndex++

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

