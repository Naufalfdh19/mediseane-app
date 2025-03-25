package converter

import (
	"montelukast/modules/order/dto"
	"montelukast/modules/order/entity"
	queryparams "montelukast/modules/order/query_params"
	productEntity "montelukast/modules/product/entity"
)

type GetUserOrdersConverter struct{}

func (c GetUserOrdersConverter) ToDto(order entity.OrderDetail) dto.GetUserOrdersResponse {
	return dto.GetUserOrdersResponse{
		OrderID: order.ID,
		Status: order.Status,
		CreatedAt: order.CreatedAt,
	}
}

type GetUserProductOrdersConverter struct{}

func (c GetUserProductOrdersConverter) ToDto(order entity.OrderProductDetail) dto.GetUserOrderDetailsResponse {
	return dto.GetUserOrderDetailsResponse{
		OrderID: order.ID,
		Status: order.Status,
		CreatedAt: order.CreatedAt,	
	}
}


type ProductOrdersConverter struct {}

func (c ProductOrdersConverter) ToDto(productOrder productEntity.ProductDetail ) dto.GetUserProductOrdersResponse {
	return dto.GetUserProductOrdersResponse{
		ID: productOrder.ID,
		Name: productOrder.Name,
		Quantity: productOrder.Stock,
		Image: productOrder.Image,
	}
}


type PaginationConverter struct{}

func (c PaginationConverter) ToDto(pagination entity.Pagination) dto.Pagination {
	return dto.Pagination{
		CurrentPage:  pagination.CurrentPage,
		TotalOrder: pagination.TotalOrder,
		TotalPage:    pagination.TotalPage,
	}
}

type QueryParamsConverter struct{}

func (c QueryParamsConverter) ToEntity(queryParams queryparams.QueryParamsDto) queryparams.QueryParams {
	return queryparams.QueryParams{
		Status:   queryParams.Status,
		SortBy: queryParams.SortBy,
		Order:  queryParams.Order,
		Limit:  queryParams.Limit,
		Page:   queryParams.Page,
	}
}
