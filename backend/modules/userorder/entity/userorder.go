package entity

import (
	"mime/multipart"
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID           int
	UserID       int
	TotalPrice   decimal.Decimal
	CreatedAt    time.Time
	OrderDetails []OrderDetail
}

type OrderDetail struct {
	ID                  int
	OrderID             int
	PharmacyID          int
	PharmacyName        string
	Status              string
	LogisticPrice       decimal.Decimal
	OrderProductDetails []OrderProductDetail
	Order
}

type OrderProductDetail struct {
	ID                int
	PharmacyProductID int
	Name              string
	Manufacturer      string
	Image             string
	Quantity          int
	Subtotal          decimal.Decimal
	OrderDetail
}


type File struct {
	File multipart.File
}