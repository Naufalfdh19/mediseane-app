package handler

import (
	"montelukast/modules/order/converter"
	"montelukast/modules/order/dto"
	queryparams "montelukast/modules/order/query_params"
	"montelukast/modules/order/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	u usecase.OrderUsecase
}

func NewOrderHandler(u usecase.OrderUsecase) OrderHandler {
	return OrderHandler{
		u: u,
	}
}

func (h *OrderHandler) GetOrdersHandler(c *gin.Context) {
	rawUserID, isExists := c.Get("user_id")
	if !isExists {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, apperror.ErrTokenInvalid)
		c.Error(err)
		return
	}
	userID, err := strconv.Atoi(rawUserID.(string))
	if err != nil {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, err)
		c.Error(err)
		return
	}

	queryParamsDto := queryparams.QueryParamsDto{}
	if err := c.ShouldBindQuery(&queryParamsDto); err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetUserOrders, apperror.ErrInvalidJSON, apperror.ErrInvalidJSON)
		c.Error(err)
		return
	}

	orders, err := h.u.GetOrders(c, converter.QueryParamsConverter{}.ToEntity(queryParamsDto), userID)
	if err != nil {
		c.Error(err)
		return
	}

	ordersDto := []dto.GetUserOrdersResponse{}

	for _, order := range orders.Orders {
		ordersDto = append(ordersDto, converter.GetUserOrdersConverter{}.ToDto(order))
	}

	var ordersList dto.OrdersList
	ordersList.Orders = ordersDto

	ordersList.Pagination = converter.PaginationConverter{}.ToDto(orders.Pagination)

	response := wrapper.ResponseData(ordersList, "get user orders success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) GetOrderedProductsHandler(c *gin.Context) {
	rawUserID, isExists := c.Get("user_id")
	if !isExists {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, apperror.ErrTokenInvalid)
		c.Error(err)
		return
	}
	userID, err := strconv.Atoi(rawUserID.(string))
	if err != nil {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, err)
		c.Error(err)
		return
	}

	rawOrderDetailsID := c.Param("id")
	orderDetailsID, err := strconv.Atoi(rawOrderDetailsID)
	if err != nil {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, err)
		c.Error(err)
		return
	}

	productOrders, err := h.u.GetOrderedProducts(c, orderDetailsID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	productOrdersDto := converter.GetUserProductOrdersConverter{}.ToDto(*productOrders)


	for _, order := range productOrders.ProductDetails {
		productOrdersDto.ProductDetails = append(productOrdersDto.ProductDetails, converter.ProductOrdersConverter{}.ToDto(order))
	}


	response := wrapper.ResponseData(productOrdersDto, "get user product orders success", nil)
	c.JSON(http.StatusOK, response)
}

func (h OrderHandler) DeleteOrderHandler(c *gin.Context) {
	rawUserID, isExists := c.Get("user_id")
	if !isExists {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, apperror.ErrTokenInvalid)
		c.Error(err)
		return
	}
	userID, err := strconv.Atoi(rawUserID.(string))
	if err != nil {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, err)
		c.Error(err)
		return
	}

	rawOrderDetailID := c.Param("id")
	orderDetailID, err := strconv.Atoi(rawOrderDetailID)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrDeleteOrder, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return 
	}

	err = h.u.DeleteOrder(c, orderDetailID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "delete order success", nil)
	c.JSON(http.StatusOK, response)

}

func (h OrderHandler) UpdateOrderStatusHandler(c *gin.Context) {
	rawUserID, isExists := c.Get("user_id")
	if !isExists {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, apperror.ErrTokenInvalid)
		c.Error(err)
		return
	}
	userID, err := strconv.Atoi(rawUserID.(string))
	if err != nil {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, err)
		c.Error(err)
		return
	}

	rawOrderDetailID := c.Param("id")
	orderDetailID, err := strconv.Atoi(rawOrderDetailID)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrDeleteOrder, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return 
	}

	err = h.u.UpdateOrderStatus(c, orderDetailID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "update order status success", nil)
	c.JSON(http.StatusOK, response)
}