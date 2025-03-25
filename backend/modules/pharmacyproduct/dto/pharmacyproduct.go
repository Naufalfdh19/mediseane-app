package dto

import (
	productDto "montelukast/modules/product/dto"
	"time"

	"github.com/shopspring/decimal"
)

type AddPharmacyProductRequest struct {
	ProductID int             `json:"product_id" binding:"required"`
	Stock     int             `json:"stock"`
	Price     decimal.Decimal `json:"price" binding:"required"`
	IsActive  bool            `json:"is_active"`
}

type UpdatePharmacyProductRequest struct {
	Stock    int  `json:"stock" binding:"required,gte=1"`
	IsActive bool `json:"is_active"`
}

type GetPharmacyProductResponse struct {
	ID                    int       `json:"id"`
	Name                  string    `json:"name"`
	GenericName           string    `json:"generic_name"`
	Manufacturer          string    `json:"manufacturer"`
	ProductClassification string    `json:"product_classification"`
	ProductForm           string    `json:"product_form"`
	Stock                 int       `json:"stock"`
	Price                 string    `json:"price"`
	IsActive              bool      `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
}
type ProductsList struct {
	Pagination productDto.Pagination        `json:"pagination"`
	Products   []GetPharmacyProductResponse `json:"pharmacy_products"`
}
