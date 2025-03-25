package dto

import (
	"mime/multipart"

	"github.com/shopspring/decimal"
)

type GetProductsResponse struct {
	ID                int             `json:"id"`
	PharmacyProductID int             `json:"pharmacy_product_id"`
	Image             string          `json:"image"`
	Name              string          `json:"name"`
	Manufacture       string          `json:"manufacture"`
	PharmacyName      string          `json:"pharmacy_name"`
	Price             decimal.Decimal `json:"price"`
}

type GetProductDetailResponse struct {
	ID                int             `json:"id"`
	PharmacyProductID int             `json:"pharmacy_product_id"`
	ProductCategories []string        `json:"product_categories"`
	Name              string          `json:"name"`
	GenericName       string          `json:"generic_name"`
	Manufacture       string          `json:"manufacture"`
	Description       string          `json:"description"`
	Image             string          `json:"image"`
	UnitInPack        int             `json:"unit_in_pack"`
	Stock             int             `json:"stock"`
	Price             decimal.Decimal `json:"price"`
	PharmacyAddress   string          `json:"address"`
	PharmacysName     string          `json:"pharmacies_name"`
}

type ProductResponse struct {
	PharmacyProductID       int     `json:"pharmacy_product_id"`
	ProductCategoriesID     []int   `json:"product_categories_id"`
	ProductClassificationID int     `json:"product_classification_id"`
	ProductFormID           *int    `json:"product_form_id"`
	Name                    string  `json:"name"`
	GenericName             string  `json:"generic_name"`
	Manufacture             string  `json:"manufacture"`
	Description             string  `json:"description"`
	Image                   string  `json:"image"`
	UnitInPack              *int    `json:"unit_in_pack"`
	Weight                  float64 `json:"weight"`
	Height                  float64 `json:"height"`
	Length                  float64 `json:"length"`
	Width                   float64 `json:"width"`
}

type AddProductRequest struct {
	ProductCategoriesID     []int   `json:"product_categories_id" binding:"required"`
	ProductClassificationID int     `json:"product_classification_id" binding:"required,gte=1"`
	ProductFormID           *int    `json:"product_form_id"`
	Name                    string  `json:"name" binding:"required,max=75,min=4"`
	GenericName             string  `json:"generic_name" binding:"required"`
	Manufacture             string  `json:"manufacture" binding:"required"`
	Description             string  `json:"description" binding:"required"`
	UnitInPack              *int    `json:"unit_in_pack"`
	Weight                  float64 `json:"weight" binding:"required,gte=0"`
	Height                  float64 `json:"height" binding:"required,gte=0"`
	Length                  float64 `json:"length" binding:"required,gte=0"`
	Width                   float64 `json:"width" binding:"required,gte=0"`
	IsActive                bool    `json:"is_active" binding:"required"`
}

type ProductDetail struct {
	ProductResponse
	Stock           int     `json:"stock"`
	Price           float64 `json:"price"`
	PharmacyAddress string  `json:"address"`
	PharmacysName   string  `json:"pharmacies_name"`
}

type GetProductDetailRequest struct {
	PharmacistProductID int `json:"pharmacists_product_id" binding:"required"`
}

type GetProductResponse struct {
	ID                    int     `json:"id"`
	ProductClassification string  `json:"product_classification"`
	ProductForm           *string `json:"product_form,omitempty"`
	Name                  string  `json:"name"`
	GenericName           string  `json:"generic_name"`
	Manufacture           string  `json:"manufacture"`
	IsActive              bool    `json:"is_active"`
}

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPage    int `json:"total_page"`
	TotalProduct int `json:"total_product"`
}

type ProductsList struct {
	Pagination Pagination            `json:"pagination"`
	Products   []GetProductsResponse `json:"products"`
}

type ProductsListAdmin struct {
	Pagination Pagination           `json:"pagination"`
	Products   []GetProductResponse `json:"products"`
}

type GetMasterProductResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MasterProductList struct {
	Pagination Pagination        `json:"pagination"`
	Products   []GetMasterProductResponse `json:"products"`
}

type FileRequest struct {
	File multipart.File `json:"file,omitempty"`
}
