package entity

import "montelukast/pkg/pagination"

type Category struct {
	ID        int
	Name      string
	UpdatedAt string
}

type PaginatedCategories struct {
	Pagination pagination.Pagination
	Categories []Category
}
