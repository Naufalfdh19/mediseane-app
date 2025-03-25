package converter

import (
	"montelukast/modules/cart/dto"
	"montelukast/modules/cart/entity"
)

type AddToCartConverter struct{}

func (c AddToCartConverter) ToEntity(cartReq dto.AddToCartRequest) entity.CartItem {
	return entity.CartItem{
		UserID:            cartReq.UserID,
		PharmacyProductID: cartReq.PharmacyProductID,
		Quantity:          cartReq.Quantity,
	}
}

type GetGroupedCartItemsConverter struct{}

func (c GetGroupedCartItemsConverter) ToDto(groupedCartItem entity.GroupedCartItem) dto.GroupedCartItemResponse {
	total := len(groupedCartItem.Items)
	cartItemResponses := make([]dto.CartItemResponse, total)

	for i := 0; i < total; i++ {
		item := groupedCartItem.Items[i]
		cartItemResponse := dto.CartItemResponse{
			CartItemID:        item.ID,
			PharmacyProductID: item.PharmacyProductID,
			Name:              item.Name,
			Manufacturer:      item.Manufacturer,
			Image:             item.Image,
			Quantity:          item.Quantity,
			Subtotal:          item.Subtotal.String(),
		}
		cartItemResponses[i] = cartItemResponse
	}

	return dto.GroupedCartItemResponse{
		PharmacyID:   groupedCartItem.PharmacyID,
		PharmacyName: groupedCartItem.PharmacyName,
		Items:        cartItemResponses,
	}
}

type GetCartItemsConverter struct{}

func (c GetCartItemsConverter) ToDto(cartItem entity.CartItem) dto.CartItemResponse {
	return dto.CartItemResponse{
		CartItemID:        cartItem.ID,
		PharmacyProductID: cartItem.PharmacyProductID,
		Quantity:          cartItem.Quantity,
	}
}
