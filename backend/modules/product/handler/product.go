package handler

import (
	"montelukast/modules/product/converter"
	"montelukast/modules/product/dto"
	queryparams "montelukast/modules/product/queryparams"
	"montelukast/modules/product/usecase"
	"strings"

	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	u usecase.ProductUsecase
}

func NewProductHandler(u usecase.ProductUsecase) ProductHandler {
	return ProductHandler{
		u: u,
	}
}

func (h *ProductHandler) GetUserProductsHandler(c *gin.Context) {
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
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetUserProduct, apperror.ErrInvalidJSON, apperror.ErrInvalidJSON)
		c.Error(err)
		return
	}

	products, err := h.u.GetUserProducts(c, converter.QueryParamsConverter{}.ToEntity(queryParamsDto), userID)
	if err != nil {
		c.Error(err)
		return
	}

	productsDto := []dto.GetProductsResponse{}

	for _, product := range products.Products {
		productsDto = append(productsDto, converter.GetUserProductsConverter{}.ToDto(product))
	}

	var productsList dto.ProductsList
	productsList.Products = productsDto

	productsList.Pagination = converter.PaginationConverter{}.ToDto(products.Pagination)

	response := wrapper.ResponseData(productsList, "get user products success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetUserProductsHomepageHandler(c *gin.Context) {
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
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetUserProduct, apperror.ErrInvalidJSON, apperror.ErrInvalidJSON)
		c.Error(err)
		return
	}

	products, err := h.u.GetUserProductsHomePage(c, converter.QueryParamsConverter{}.ToEntity(queryParamsDto), userID)
	if err != nil {
		c.Error(err)
		return
	}

	productsDto := []dto.GetProductsResponse{}

	for _, product := range products.Products {
		productsDto = append(productsDto, converter.GetUserProductsConverter{}.ToDto(product))
	}

	var productsList dto.ProductsList
	productsList.Products = productsDto

	productsList.Pagination = converter.PaginationConverter{}.ToDto(products.Pagination)

	response := wrapper.ResponseData(productsList, "get user products success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetProductDetailHandler(c *gin.Context) {
	pharmacyProductIDStr := c.Param("id")
	pharmacyProductID, err := strconv.Atoi(pharmacyProductIDStr)
	if err != nil {
		c.Error(apperror.NewErrStatusBadRequest(appconstant.FieldErrGetProduct, apperror.ErrQueryParams, err))
		return
	}

	res, err := h.u.GetProductDetail(c, pharmacyProductID)
	if err != nil {
		c.Error(err)
		return
	}

	resDto := converter.ProductDetailConverter{}.ToDto(*res)

	response := wrapper.ResponseData(resDto, "get product detail success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetMasterProductsHandler(c *gin.Context) {
	queryParamsDto := queryparams.QueryParamsDto{}
	if err := c.ShouldBindQuery(&queryParamsDto); err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrGetPharmacyProducts, apperror.ErrInvalidJSON, apperror.ErrInvalidJSON)
		c.Error(err)
		return
	}

	products, err := h.u.GetMasterProducts(c, converter.QueryParamsConverter{}.ToEntity(queryParamsDto))
	if err != nil {
		c.Error(err)
		return
	}

	productsDto := []dto.GetMasterProductResponse{}

	for _, product := range products.Products {
		productsDto = append(productsDto, converter.GetMasterProductConverter{}.ToDto(product))
	}

	var productsList dto.MasterProductList
	productsList.Products = productsDto

	productsList.Pagination = converter.PaginationConverter{}.ToDto(products.Pagination)

	response := wrapper.ResponseData(productsList, "get products success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) AddProductHandler(c *gin.Context) {
	productReq := dto.AddProductRequest{}

	err := apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrAddProduct, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}

	err = c.ShouldBindJSON(&productReq)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.u.AddProduct(c, converter.AddProductsConverter{}.ToEntity(productReq))
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "add product success!", nil)
	c.JSON(http.StatusCreated, response)
}

func (h *ProductHandler) UpdateProductHandler(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdateProduct, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	err = apperror.JsonValidator(c)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdateProduct, apperror.ErrInvalidJSON, err)
		c.Error(err)
		return
	}

	productReq := dto.AddProductRequest{}

	err = c.ShouldBindJSON(&productReq)
	if err != nil {
		c.Error(err)
		return
	}

	productEntity := converter.AddProductsConverter{}.ToEntity(productReq)
	productEntity.ID = productID

	err = h.u.UpdateProduct(c, productEntity)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "update product success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) UpdateProductPhotoHandler(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdateProduct, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	_, fileHeader, err := c.Request.FormFile("file")
	fileName := fileHeader.Filename
	extension := fileName[strings.Index(fileName, ".")+1:]
	if err != nil {
		response := wrapper.ResponseData(nil, "Select a file to upload", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if extension != "jpeg" && extension != "png" && extension != "jpg" {
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

	url, err := h.u.UpdateProductPhoto(c, converter.FileConverter{}.ToEntity(fileRequest), productID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(url, "update product photo success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) DeleteProductHandler(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldErrDeleteProduct, apperror.ErrConvertVariableType, err)
		c.Error(err)
		return
	}

	err = h.u.DeleteProduct(c, productID)
	if err != nil {
		c.Error(err)
		return
	}

	response := wrapper.ResponseData(nil, "delete product success!", nil)
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetProductsAdminHandler(c *gin.Context) {
	queryParamsDto := queryparams.AdminQueryParamsDto{}
	if err := c.ShouldBindQuery(&queryParamsDto); err != nil {
		err := apperror.NewErrStatusBadRequest(appconstant.FieldAdminGetProducts, apperror.ErrInvalidJSON, apperror.ErrInvalidJSON)
		c.Error(err)
		return
	}

	products, err := h.u.GetProductsAdmin(c, converter.AdminQueryParamsConverter{}.ToEntity(queryParamsDto))
	if err != nil {
		c.Error(err)
		return
	}

	productsDto := []dto.GetProductResponse{}

	for _, product := range products.Products {
		productsDto = append(productsDto, converter.GetProductsAdminConverter{}.ToDto(product))
	}

	var productsList dto.ProductsListAdmin
	productsList.Products = productsDto

	productsList.Pagination = converter.PaginationConverter{}.ToDto(products.Pagination)

	response := wrapper.ResponseData(productsList, "get products success", nil)
	c.JSON(http.StatusOK, response)
}
