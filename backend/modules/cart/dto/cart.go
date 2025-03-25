package dto

type AddToCartRequest struct {
	UserID            int `json:"user_id"`
	PharmacyProductID int `json:"pharmacy_product_id" binding:"required"`
	Quantity          int `json:"quantity" binding:"required"`
}

type CartItemResponse struct {
	CartItemID        int    `json:"cart_item_id"`
	PharmacyProductID int    `json:"pharmacy_product_id"`
	Name              string `json:"name"`
	Manufacturer      string `json:"manufacturer"`
	Image             string `json:"image,omitempty"`
	Quantity          int    `json:"quantity"`
	Subtotal          string `json:"subtotal,omitempty"`
}

type GroupedCartItemResponse struct {
	PharmacyID   int                `json:"pharmacy_id"`
	PharmacyName string             `json:"pharmacy_name"`
	Items        []CartItemResponse `json:"items"`
}

type ListGroupedCartItem struct {
	ID                  string                    `json:"id"`
	ListGroupedCartItem []GroupedCartItemResponse `json:"data"`
}

type CheckoutCartRequest struct {
	IDs []int `json:"ids" binding:"required"`
}
