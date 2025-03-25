package repository

import (
	"montelukast/modules/address/entity"
)

func AddressesParam(filter entity.AddressFilter) string {
	if filter.IsActive == "true" {
		return "AND is_active = true "
	}
	return ""
}
