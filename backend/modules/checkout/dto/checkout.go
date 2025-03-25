package dto

type CheckoutData struct {
	IDCart           string         `json:"id_cart"`
	ListDeliveryData []DeliveryData `json:"delivery_data_list"`
}

type DeliveryData struct {
	PharmacyID int `json:"pharmacy_id"`
	DeliveryID int `json:"delivery_id"`
}
