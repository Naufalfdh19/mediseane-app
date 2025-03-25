package handler

import (
	"montelukast/modules/cart/converter"
	"montelukast/modules/cart/dto"
	"montelukast/modules/cart/entity"
	"montelukast/modules/cart/usecase"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	u usecase.CartUsecase
}

func NewCartHandler(u usecase.CartUsecase) CartHandler {
	return CartHandler{
		u: u,
	}
}

func (h *CartHandler) AddToCartHandler(c *gin.Context) {
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
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrAddToCart, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}

	cartReq := dto.AddToCartRequest{}
	cartReq.UserID = userID

	err = c.ShouldBindJSON(&cartReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.AddToCart(c, converter.AddToCartConverter{}.ToEntity(cartReq))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "cart updated successfully!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *CartHandler) DeleteFromCartHandler(c *gin.Context) {
	rawID := c.Param("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		err := apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		c.Error(err)
		return
	}

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

	err = h.u.DeleteFromCart(c, entity.CartItem{ID: id, UserID: userID})
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "product deleted from cart successfully!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *CartHandler) GetGroupedCartItemsHandler(c *gin.Context) {
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

	cartItems, err := h.u.GetGroupedCartItems(c, userID)
	if err != nil {
		c.Error(err)
		return
	}
	cartItemsRes := []dto.GroupedCartItemResponse{}

	var cartItemConverter converter.GetGroupedCartItemsConverter
	for _, cartItem := range cartItems {
		cartItemsRes = append(cartItemsRes, cartItemConverter.ToDto(cartItem))
	}

	response := wrapper.ResponseData(cartItemsRes, "get cart success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *CartHandler) GetCartItemsHandler(c *gin.Context) {
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

	cartItems, err := h.u.GetCartItems(c, userID)
	if err != nil {
		c.Error(err)
		return
	}
	cartItemsRes := []dto.CartItemResponse{}

	var cartItemConverter converter.GetCartItemsConverter
	for _, cartItem := range cartItems {
		cartItemsRes = append(cartItemsRes, cartItemConverter.ToDto(cartItem))
	}

	response := wrapper.ResponseData(cartItemsRes, "get cart success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *CartHandler) GetSelectedCartItemsHandler(c *gin.Context) {
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
	checkoutReq := dto.CheckoutCartRequest{}

	err = c.ShouldBindJSON(&checkoutReq)
	if err != nil {
		c.Error(err)
		return
	}

	cartItems, err := h.u.GetSelectedCartItems(c, userID, checkoutReq.IDs)
	if err != nil {
		c.Error(err)
		return
	}
	var groupedList dto.ListGroupedCartItem
	cartItemsRes := []dto.GroupedCartItemResponse{}
	var cartItemConverter converter.GetGroupedCartItemsConverter
	for _, cartItem := range cartItems.GroupedItem {
		cartItemsRes = append(cartItemsRes, cartItemConverter.ToDto(cartItem))
	}
	groupedList.ID = cartItems.ID
	groupedList.ListGroupedCartItem = cartItemsRes
	response := wrapper.ResponseData(groupedList, "get selected cart items success!", nil)
	c.JSON(http.StatusOK, response)
}
