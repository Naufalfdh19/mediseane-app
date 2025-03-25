package handler

import (
	"montelukast/modules/userorder/converter"
	dto "montelukast/modules/userorder/dto"
	entity "montelukast/modules/userorder/entity"
	"montelukast/modules/userorder/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserOrderHandler struct {
	u usecase.UserOrderUsecase
}

func NewUserOrderHandler(u usecase.UserOrderUsecase) UserOrderHandler {
	return UserOrderHandler{
		u: u,
	}
}

func (h UserOrderHandler) UpdatePaymentHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.Error(apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer))
		return
	}
	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}
	rawOrderID := c.Param("order_id")
	orderID, err := strconv.Atoi(rawOrderID)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrUploadPayment, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrUploadPayment, apperror.ErrFileEmpty, err)
		c.Error(err)
		return
	}
	fileName := fileHeader.Filename
	extension := fileName[strings.Index(fileName, ".")+1:]
	if extension != "pdf" && extension != "png" && extension != "jpg" {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrImageType, apperror.ErrUploadImage, err)
		c.Error(err)
		return
	}
	if fileHeader.Size > appconstant.IMAGESIZEMAX {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrImageSize, apperror.ErrUploadImageSize, err)
		c.Error(err)
		return
	}
	formFile, err := fileHeader.Open()
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrImageType, apperror.ErrUploadImage, err)
		c.Error(err)
		return
	}

	fileRequest := dto.FileRequest{}
	fileRequest.File = formFile

	err = h.u.UpdatePaymentStatus(c, converter.FileConverter{}.ToEntity(fileRequest), orderID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "update order status success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h UserOrderHandler) ConfirmDeliveryHandler(c *gin.Context) {
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

	rawOrderDetailID := c.Param("order-detail-id")
	orderDetailID, err := strconv.Atoi(rawOrderDetailID)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrConfirmDelivery, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	orderDetail := entity.OrderDetail{ID: orderDetailID}
	order := entity.Order{UserID: userID, OrderDetails: []entity.OrderDetail{orderDetail}}
	err = h.u.UpdateDeliveryStatus(c, order)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "update order status success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h UserOrderHandler) GetDetailedOrdersHandler(c *gin.Context) {
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

	filterDto := dto.OrderFilterRequest{UserID: userID}
	if err := c.ShouldBindQuery(&filterDto); err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetUserOrders, apperror.ErrConvertVariableType, nil)
		c.Error(err)
		return
	}

	var filterConverter converter.FilterOrdersConverter
	orders, err := h.u.GetDetailedOrders(c, filterConverter.ToEntity(filterDto))
	if err != nil {
		c.Error(err)
		return
	}

	var orderConverter converter.GetDetailedOrdersConverter
	response := wrapper.ResponseData(orderConverter.ToDto(orders), "get orders success!", nil)
	c.JSON(http.StatusOK, response)
}
