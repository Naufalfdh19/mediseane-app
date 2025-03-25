package converter

import (
	dto "montelukast/modules/userorder/dto"
	entity "montelukast/modules/userorder/entity"
)

type GetDetailedOrdersConverter struct{}

func (c GetDetailedOrdersConverter) ToDto(orders []entity.Order) []dto.OrderResponse {
	total := len(orders)
	orderResponses := make([]dto.OrderResponse, total)

	for i := 0; i < total; i++ {
		order := orders[i]
		orderResponse := dto.OrderResponse{
			ID:           order.ID,
			TotalPrice:   order.TotalPrice.String(),
			CreatedAt:    order.CreatedAt,
			OrderDetails: c.toOrderDetails(order.OrderDetails),
		}
		orderResponses[i] = orderResponse
	}

	return orderResponses
}

func (c GetDetailedOrdersConverter) toOrderDetails(orderDetails []entity.OrderDetail) []dto.OrderDetailResponse {
	total := len(orderDetails)
	orderDetailResponses := make([]dto.OrderDetailResponse, total)

	for i := 0; i < total; i++ {
		detail := orderDetails[i]
		orderDetailResponse := dto.OrderDetailResponse{
			ID:                 detail.ID,
			PharmacyID:         detail.PharmacyID,
			PharmacyName:       detail.PharmacyName,
			Status:             detail.Status,
			LogisticPrice:      detail.LogisticPrice.String(),
			OrderProductDetail: c.toOrderProductDetails(detail.OrderProductDetails),
		}
		orderDetailResponses[i] = orderDetailResponse
	}

	return orderDetailResponses
}

func (c GetDetailedOrdersConverter) toOrderProductDetails(orderProductDetail []entity.OrderProductDetail) []dto.OrderProductDetailResponse {
	total := len(orderProductDetail)
	orderProductDetailResponses := make([]dto.OrderProductDetailResponse, total)

	for i := 0; i < total; i++ {
		product := orderProductDetail[i]
		orderProductDetailResponse := dto.OrderProductDetailResponse{
			OrderProductID:    product.ID,
			PharmacyProductID: product.PharmacyProductID,
			Name:              product.Name,
			Manufacturer:      product.Manufacturer,
			Image:             product.Image,
			Quantity:          product.Quantity,
			Subtotal:          product.Subtotal.String(),
		}
		orderProductDetailResponses[i] = orderProductDetailResponse
	}

	return orderProductDetailResponses
}

type FilterOrdersConverter struct{}

func (c FilterOrdersConverter) ToEntity(filterDto dto.OrderFilterRequest) (filter entity.OrderFilter) {
	return entity.OrderFilter{
		SortBy: filterDto.SortBy,
		Order:  filterDto.Order,
		UserID: filterDto.UserID,
		Status: filterDto.Status,
	}
}


type FileConverter struct{}

func (c FileConverter) ToEntity(fileDTO dto.FileRequest) entity.File {
	return entity.File{
		File: fileDTO.File,
	}
}