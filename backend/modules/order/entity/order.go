package entity

import (
	productEntity "montelukast/modules/product/entity"
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID         int
	UserID     int
	TotalPrice decimal.Decimal
}

type OrderDetail struct {
	Order
	PharmacyID    int
	LogisticPrice decimal.Decimal
	Status        string
	CreatedAt     time.Time
}

type OrderProductDetail struct {
	OrderDetail
	Quantity int
	Price decimal.Decimal
	ProductDetails []productEntity.ProductDetail
}

type Pagination struct {
	CurrentPage int
	TotalOrder  int
	TotalPage   int
}

type OrdersList struct {
	Pagination Pagination
	Orders     []OrderDetail
}
