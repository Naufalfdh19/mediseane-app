package handler

import (
	"montelukast/modules/checkout/converter"
	"montelukast/modules/checkout/dto"
	"montelukast/modules/checkout/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CheckoutHandler struct {
	checkOutUsecase usecase.CheckoutUsecase
}

func NewCheckoutHandler(checkOutUsecase usecase.CheckoutUsecase) CheckoutHandler {
	return CheckoutHandler{
		checkOutUsecase: checkOutUsecase,
	}
}

func (h *CheckoutHandler) CheckoutCartHandler(c *gin.Context) {
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
	err = apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetCart, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}
	var converter converter.CheckoutConverterImpl
	checkoutData := dto.CheckoutData{}
	err = c.ShouldBindJSON(&checkoutData)
	if err != nil {
		c.Error(err)
		return
	}
	err = <-h.checkOutUsecase.Checkout(c, converter.ToEntity(checkoutData), userID)
	if err != nil {
		c.Error(err)
		return
	}
	response := wrapper.ResponseData(nil, "checkout success,waiting for payment", nil)
	c.JSON(http.StatusOK, response)
}
func (h CheckoutHandler) CancelOrder(c *gin.Context) {
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
	rawOrderID := c.Param("order-id")
	orderID, err := strconv.Atoi(rawOrderID)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrConfirmDelivery, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	err = h.checkOutUsecase.CancelOrderByUser(c, userID, orderID)
	if err != nil {
		c.Error(err)
		return
	}
	response := wrapper.ResponseData(nil, "cancel order success", nil)
	c.JSON(http.StatusOK, response)

}
