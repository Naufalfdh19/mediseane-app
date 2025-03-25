package repository

import (
	"fmt"
	"montelukast/modules/category/entity"
	"strings"
)

func CategoryParam(filter entity.CategoryFilter) (paramQuery string, order string, args []interface{}) {
	queryList := []string{}
	sortLimit := []string{}
	condition := []string{}
	args = []interface{}{}
	paramIndex := 1
	if filter.Name != "" {
		condition = append(condition, fmt.Sprintf("c.name ILIKE '%%' || $%d || '%%'", paramIndex))
		args = append(args, filter.Name)
		paramIndex++
	}
	if len(condition) > 0 {
		queryList = append(queryList, " AND ", strings.Join(condition, " AND "))
	}
	sortField := ""
	switch filter.SortBy {
	case "name":
		sortField = "c.name"
	default:
		sortField = "c.updated_at"
	}

	order = " ASC "
	if strings.ToLower(filter.Order) == "desc" {
		order = " DESC "
	}
	sortLimit = append(sortLimit, " ORDER BY ", sortField, order, fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1))
	query := strings.Join(queryList, "")
	sortPage := strings.Join(sortLimit, "")
	args = append(args, filter.GetLimit(), filter.GetOffset())
	return query, sortPage, args
}

func CategoryParamCount(filter entity.CategoryFilterCount) (paramQuery string, args []interface{}) {
	queryList := []string{}
	condition := []string{}
	args = []interface{}{}
	paramIndex := 1
	if filter.Name != "" {
		condition = append(condition, fmt.Sprintf("c.name ILIKE '%%' || $%d || '%%'", paramIndex))
		args = append(args, filter.Name)
		paramIndex++
	}
	if len(condition) > 0 {
		queryList = append(queryList, " AND ", strings.Join(condition, " AND "))
	}
	query := strings.Join(queryList, "")
	return query, args
}
