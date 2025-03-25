package dto

import (
	"time"
)

type GetUserOrdersResponse struct {
	OrderID   int       `json:"order_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type GetUserOrderDetailsResponse struct {
	OrderID        int                           `json:"order_id"`
	Status         string                        `json:"status"`
	CreatedAt      time.Time                     `json:"created_at"`
	ProductDetails []GetUserProductOrdersResponse `json:"product_list"`
}

type GetUserProductOrdersResponse struct {
	ID       int    `json:"product_id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Image    string `json:"image"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	TotalOrder  int `json:"total_order"`
}

type OrdersList struct {
	Pagination Pagination              `json:"pagination"`
	Orders     []GetUserOrdersResponse `json:"orders"`
}
