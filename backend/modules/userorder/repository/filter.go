package repository

import (
	"fmt"
	"montelukast/modules/userorder/entity"
	"strings"
)

func OrderParam(filter entity.OrderFilter) (paramQuery string, order string, args []interface{}) {
	queryList := []string{}
	sortLimit := []string{}
	condition := []string{}
	args = []interface{}{}
	paramIndex := 1
	
	if filter.Status != "" {
		condition = append(condition, fmt.Sprintf("od.status = $%d", paramIndex))
		args = append(args, filter.Status)
		paramIndex++
	}
	if len(condition) > 0 {
		queryList = append(queryList, " AND ", strings.Join(condition, " AND "))
	}

	sortField := "o.created_at"

	order = " DESC "
	if strings.ToLower(filter.Order) == "asc" {
		order = " ASC "
	}

	sortLimit = append(sortLimit, " ORDER BY ", sortField, order)
	query := strings.Join(queryList, "")
	sortPage := strings.Join(sortLimit, "")
	return query, sortPage, args
}
