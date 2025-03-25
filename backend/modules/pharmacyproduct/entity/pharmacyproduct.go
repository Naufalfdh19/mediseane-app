package entity

import "github.com/shopspring/decimal"

type PharmacyProduct struct {
	ID int
	PharmacyID   int
	PharmacyName string
	ProductID    int
	ProductName  string
	Stock        int
	Price        decimal.Decimal
	IsActive     bool
}


