package repository

import (
	"fmt"
	"montelukast/modules/pharmacy/entity"
	"strings"
)

func PharmacyParam(filter entity.PharmacyFilter) (paramQuery string, order string, args []interface{}) {
	queryList := []string{}
	sortLimit := []string{}
	condition := []string{}
	args = []interface{}{}
	paramIndex := 1
	if filter.City != "" {
		condition = append(condition, fmt.Sprintf("p.city ILIKE '%%' || $%d || '%%'", paramIndex))
		args = append(args, filter.City)
		paramIndex++
	}
	if filter.Name != "" {
		condition = append(condition, fmt.Sprintf("p.name ILIKE '%%' || $%d || '%%'", paramIndex))
		args = append(args, filter.Name)
		paramIndex++
	}
	if len(condition) > 0 {
		queryList = append(queryList, " AND ", strings.Join(condition, " AND "))
	}
	sortField := ""
	switch filter.Field {
	case "last_update":
		sortField = "p.updated_at"
	default:
		sortField = "p.updated_at"
	}

	order = " DESC "
	if strings.ToLower(filter.Order) == "asc" {
		order = " ASC "
	}
	sortLimit = append(sortLimit, " ORDER BY ", sortField, order, fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1))
	query := strings.Join(queryList, "")
	sortPage := strings.Join(sortLimit, "")
	args = append(args, filter.GetLimit(), filter.GetOffset())
	return query, sortPage, args

}

func PharmacyParamCount(filter entity.PharmacyFilterCount) (paramQuery string, args []interface{}) {
	queryList := []string{}
	condition := []string{}
	args = []interface{}{}
	paramIndex := 1
	if filter.City != "" {
		condition = append(condition, fmt.Sprintf("p.city ILIKE '%%' || $%d || '%%'", paramIndex))
		args = append(args, filter.City)
		paramIndex++
	}
	if filter.Name != "" {
		condition = append(condition, fmt.Sprintf("p.name ILIKE '%%' || $%d || '%%'", paramIndex))
		args = append(args, filter.Name)
		paramIndex++
	}
	if len(condition) > 0 {
		queryList = append(queryList, " AND ", strings.Join(condition, " AND "))
	}
	query := strings.Join(queryList, "")
	return query, args
}
