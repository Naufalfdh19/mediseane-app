package entity

import "github.com/shopspring/decimal"

type GroupedCartItem struct {
	PharmacyID   int
	PharmacyName string
	Items        []CartItem
}

type ListGroupedCartItem struct {
	ID          string
	GroupedItem []GroupedCartItem
}

type GroupedCheckoutItem struct {
	PharmacyID   int
	PharmacyName string
	Items        []CartItem
	Logistic     []LogisticPartner
}

type CheckoutData struct {
	IDCart           string
	ListDeliveryData []DeliveryData
}

type DeliveryData struct {
	PharmacyID int
	DeliveryID int
}

type LogisticPartner struct {
	LogisticPartnerID int
	LogisticName      string
	Etd               int
}

type DeliveryPriceData struct {
	PharmacyID    int
	LogisticPrice decimal.Decimal
	Status        string
}

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
