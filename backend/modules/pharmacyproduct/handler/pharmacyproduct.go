package handler

import (
	"montelukast/modules/pharmacyproduct/converter"
	"montelukast/modules/pharmacyproduct/dto"
	"montelukast/modules/pharmacyproduct/entity"
	"montelukast/modules/pharmacyproduct/queryparams"
	"montelukast/modules/pharmacyproduct/usecase"
	productConverter "montelukast/modules/product/converter"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PharmacyProductHandler struct {
	u usecase.PharmacyProductUsecase
}

func NewPharmacyProductHandler(u usecase.PharmacyProductUsecase) PharmacyProductHandler {
	return PharmacyProductHandler{
		u: u,
	}
}

func (h PharmacyProductHandler) AddPharmacyProductHandler(c *gin.Context) {
	rawPharmacistID, exists := c.Get("user_id")
	if !exists {
		c.Error(apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer))
		return
	}
	pharmacistID, err := strconv.Atoi(rawPharmacistID.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	err = apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPharmacyProduct, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}

	pharmacyProductReq := dto.AddPharmacyProductRequest{}
	err = c.ShouldBindJSON(&pharmacyProductReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.AddPharmacyProduct(c, converter.AddPharmacyProductConverter{}.ToEntity(pharmacyProductReq), pharmacistID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "add pharmacy product success!", nil)
	c.JSON(http.StatusCreated, response)
}

func (h PharmacyProductHandler) UpdatePharmacyProductHandler(c *gin.Context) {
	rawPharmacistID, exists := c.Get("user_id")
	if !exists {
		c.Error(apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer))
		return
	}
	pharmacistID, err := strconv.Atoi(rawPharmacistID.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	rawPharmacyProductID := c.Param("id")
	pharmacyProductID, err := strconv.Atoi(rawPharmacyProductID)
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	err = apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPharmacyProduct, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}

	pharmacyProductReq := dto.UpdatePharmacyProductRequest{}
	err = c.ShouldBindJSON(&pharmacyProductReq)
	if err != nil {
		c.Error(err)
		return
	}

	pharmacistProductEnt := converter.UpdatePharmacyProductConverter{}.ToEntity(pharmacyProductReq)
	pharmacistProductEnt.ID = pharmacyProductID
	err = h.u.UpdatePharmacyProduct(c, pharmacistProductEnt, pharmacistID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "update pharmacy product success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h PharmacyProductHandler) DeletePharmacyProductHandler(c *gin.Context) {
	rawPharmacistID, exists := c.Get("user_id")
	if !exists {
		c.Error(apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer))
		return
	}
	pharmacistID, err := strconv.Atoi(rawPharmacistID.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	rawPharmacyProductID := c.Param("id")
	pharmacyProductID, err := strconv.Atoi(rawPharmacyProductID)
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	pharmacistProductEnt := entity.PharmacyProduct{}
	pharmacistProductEnt.ID = pharmacyProductID
	err = h.u.DeletePharmacyProduct(c, pharmacistProductEnt, pharmacistID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "delete pharmacy product success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h PharmacyProductHandler) GetPharmacyProductDetailHandler(c *gin.Context) {
	rawPharmacistID, exists := c.Get("user_id")
	if !exists {
		c.Error(apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer))
		return
	}
	pharmacistID, err := strconv.Atoi(rawPharmacistID.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	rawID := c.Param("id")
	pharmacyProductID, err := strconv.Atoi(rawID)
	if err != nil {
		err = apperror.NewErrStatusBadRequest(appconstant.FieldErrGetProductDetail, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	pharmacyProduct, err := h.u.GetPharmacyProductDetail(c, pharmacyProductID, pharmacistID)
	if err != nil {
		c.Error(err)
		return
	}

	pharmacyProductDto := converter.GetPharmacyProductConverter{}.ToDto(*pharmacyProduct)

	response := wrapper.ResponseData(pharmacyProductDto, "get pharmacy product success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *PharmacyProductHandler) GetPharmacyProductsHandler(c *gin.Context) {
	rawPharmacistID, exists := c.Get("user_id")
	if !exists {
		c.Error(apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, apperror.ErrInternalServer))
		return
	}
	pharmacistID, err := strconv.Atoi(rawPharmacistID.(string))
	if err != nil {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err))
		return
	}

	queryParamsDto := queryparams.QueryParamsDto{}
	if err := c.ShouldBindQuery(&queryParamsDto); err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetPharmacyProducts, apperror.ErrInvalidJSON, apperror.ErrInvalidJSON)
		c.Error(err)
		return
	}

	products, err := h.u.GetPharmacyProducts(c, converter.QueryParamsConverter{}.ToEntity(queryParamsDto), pharmacistID)
	if err != nil {
		c.Error(err)
		return
	}

	productsDto := []dto.GetPharmacyProductResponse{}

	for _, product := range products.Products {
		productsDto = append(productsDto, converter.GetPharmacyProductConverter{}.ToDto(product))
	}

	var productsList dto.ProductsList
	productsList.Products = productsDto

	productsList.Pagination = productConverter.PaginationConverter{}.ToDto(products.Pagination)

	response := wrapper.ResponseData(productsList, "get products success", nil)
	c.JSON(http.StatusOK, response)
}
