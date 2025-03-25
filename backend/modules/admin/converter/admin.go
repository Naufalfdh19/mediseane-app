package converter

import (
	"montelukast/modules/admin/dto"
	"montelukast/modules/admin/entity"
)


type PaginationConverter struct{}

func (c PaginationConverter) ToDto(pagination entity.Pagination) dto.Pagination {
	return dto.Pagination{
		CurrentPage: pagination.CurrentPage,
		TotalPage: pagination.TotalPage,
		TotalUser: pagination.TotalUser,
	}
}

type GetusersConverter struct{}

func (c GetusersConverter) ToDto(user entity.User) dto.GetusersResponse {
	return dto.GetusersResponse{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		ProfilePhoto: user.ProfilePhoto,
		Role: user.Role,
	}
}