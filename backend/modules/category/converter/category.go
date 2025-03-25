package converter

import (
	"montelukast/modules/category/dto"
	"montelukast/modules/category/entity"
)

type CategoryAddConverter struct{}

func (c CategoryAddConverter) ToEntity(categoryReq dto.CategoryAddRequest) entity.Category {
	return entity.Category{
		Name: categoryReq.Name,
	}
}

type CategoryUpdateConverter struct{}

func (c CategoryUpdateConverter) ToEntity(categoryReq dto.CategoryUpdateRequest) entity.Category {
	return entity.Category{
		ID:   categoryReq.ID,
		Name: categoryReq.Name,
	}
}

type CategoryDeleteConverter struct{}

func (c CategoryDeleteConverter) ToEntity(categoryReq dto.CategoryDeleteRequest) entity.Category {
	return entity.Category{
		ID: categoryReq.ID,
	}
}

type GetCategoriesConverter struct{}

func (c GetCategoriesConverter) ToDto(category entity.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		UpdatedAt: category.UpdatedAt,
	}
}

type FilterCategoriesConverter struct{}

func (c FilterCategoriesConverter) ToEntity(filterDTO dto.CategoryFilterRequest) (filter entity.CategoryFilter) {
	return entity.CategoryFilter{
		SortBy: filterDTO.SortBy,
		Order:  filterDTO.Order,
		Name:   filterDTO.Name,
		Limit:  filterDTO.Limit,
		Page:   filterDTO.Page,
	}
}
