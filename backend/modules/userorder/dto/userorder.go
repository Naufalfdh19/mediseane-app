package entity

import (
	"mime/multipart"
	"time"
)

type OrderResponse struct {
	ID           int                   `json:"id"`
	TotalPrice   string                `json:"total_price"`
	CreatedAt    time.Time             `json:"order_date"`
	OrderDetails []OrderDetailResponse `json:"order_details"`
}

type OrderDetailResponse struct {
	ID                 int                          `json:"details_id"`
	PharmacyID         int                          `json:"pharmacy_id"`
	PharmacyName       string                       `json:"pharmacy_name"`
	Status             string                       `json:"status"`
	LogisticPrice      string                       `json:"logistic_price"`
	OrderProductDetail []OrderProductDetailResponse `json:"order_products"`
}

type OrderProductDetailResponse struct {
	OrderProductID    int    `json:"order_product_id"`
	PharmacyProductID int    `json:"pharmacy_product_id"`
	Name              string `json:"name"`
	Manufacturer      string `json:"manufacturer"`
	Image             string `json:"image"`
	Quantity          int    `json:"quantity"`
	Subtotal          string `json:"subtotal"`
}

type OrderFilterRequest struct {
	SortBy string `form:"sortBy"`
	Order  string `form:"order"`
	UserID int    `form:"user_id"`
	Status string `form:"status"`
}

type FileRequest struct {
	File multipart.File `json:"file,omitempty"`
}