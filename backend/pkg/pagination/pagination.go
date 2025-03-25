package pagination

type Pagination struct {
	CurrentPage int
	TotalItem   int
	TotalPage   int
	Limit       int
}

type PaginationResponse struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	TotalItem   int `json:"total_item"`
}

type PaginationConverter struct{}

func (c PaginationConverter) ToDto(pagination Pagination) PaginationResponse {
	return PaginationResponse{
		CurrentPage: pagination.CurrentPage,
		TotalItem:   pagination.TotalItem,
		TotalPage:   pagination.TotalPage,
	}
}
