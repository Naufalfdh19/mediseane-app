package dto

import "montelukast/pkg/pagination"

type CategoryAddRequest struct {
	Name string `json:"name" binding:"required,min=3"`
}

type CategoryUpdateRequest struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type CategoryDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

type CategoryResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at"`
}

type PaginatedCategoriesResponse struct {
	Pagination pagination.PaginationResponse `json:"pagination"`
	Categories []CategoryResponse            `json:"list_item"`
}

type CategoryFilterRequest struct {
	SortBy string `form:"sortBy"`
	Order  string `form:"order"`
	Name   string `form:"filter"`
	Limit  int    `form:"limit"`
	Page   int    `form:"page"`
}
