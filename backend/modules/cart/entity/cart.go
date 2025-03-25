package entity

import "github.com/shopspring/decimal"

type CartItem struct {
	ID                int
	UserID            int
	PharmacyProductID int
	Name              string
	Manufacturer      string
	Image             string
	Quantity          int
	Subtotal          decimal.Decimal
	PharmacyID        int
	PharmacyName      string
}

type GroupedCartItem struct {
	ID           string
	PharmacyID   int
	PharmacyName string
	Items        []CartItem
}

type ListGroupedCartItem struct {
	ID          string
	GroupedItem []GroupedCartItem
}
