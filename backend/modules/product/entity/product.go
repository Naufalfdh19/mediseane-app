package entity

import (
	"mime/multipart"

	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID                      int
	PharmacyProductID       int
	PharmacistProductID     int
	ProductCategories       []string
	ProductCategoriesID     []int
	ProductClassification   string
	ProductClassificationID int
	ProductForm             string
	ProductFormID           *int
	Name                    string
	GenericName             string
	Manufacture             string
	Description             string
	Image                   string
	UnitInPack              *int
	Weight                  float64
	Height                  float64
	Length                  float64
	Width                   float64
	IsActive                bool
}

type ProductAdmin struct {
	ID                    int
	PharmacistProduct     string
	ProductCategory       string
	ProductClassification string
	ProductForm           *string
	Name                  string
	GenericName           string
	Manufacture           string
	Description           string
	Image                 string
	UnitInPack            *int
	Weight                float64
	Height                float64
	Length                float64
	Width                 float64
	IsActive bool
}

type ProductDetail struct {
	Product
	Stock           int
	Price           decimal.Decimal
	PharmacyAddress string
	PharmacyName    string
	CreatedAt       time.Time
}

type Pagination struct {
	CurrentPage  int
	TotalPage    int
	TotalProduct int
}

type ProductsList struct {
	Pagination Pagination
	Products   []ProductDetail
}

type ProductListAdmin struct {
	Pagination Pagination
	Products   []ProductAdmin
}

type File struct {
	File multipart.File
}

type CategoryBoundary struct {
	Minimum int
	Maximum int
}