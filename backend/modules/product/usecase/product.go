package usecase

import (
	"context"
	"math"
	"montelukast/modules/product/entity"
	queryparams "montelukast/modules/product/queryparams"
	"montelukast/modules/product/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"

	"montelukast/pkg/transaction"
)

type ProductUsecase interface {
	GetUserProducts(c context.Context, queryParams queryparams.QueryParams, userID int) (*entity.ProductsList, error)
	GetGeneralProducts(c context.Context, queryParams queryparams.QueryParams) (*entity.ProductsList, error)
	GetUserProductsHomePage(c context.Context, queryParams queryparams.QueryParams, userID int) (*entity.ProductsList, error)
	GetGeneralProductsHomePage(c context.Context, queryParams queryparams.QueryParams) (*entity.ProductsList, error)
	UpdateProduct(c context.Context, product entity.Product) error
	GetProductsAdmin(c context.Context, queryParams queryparams.AdminQueryParams) (*entity.ProductListAdmin, error)
	DeleteProduct(c context.Context, productID int) error
	GetProductDetail(c context.Context, pharmacistsProductID int) (*entity.ProductDetail, error)
	GetMasterProducts(c context.Context, queryParams queryparams.QueryParams) (*entity.ProductsList, error)
	AddProduct(c context.Context, product entity.Product) error
	UpdateProductPhoto(c context.Context, file entity.File, productID int) (string, error) 
}

type productUsecaseImpl struct {
	r  repository.ProductRepo
	tr transaction.TransactorRepoImpl
}

func NewProductUsecase(r repository.ProductRepo, tr transaction.TransactorRepoImpl) productUsecaseImpl {
	return productUsecaseImpl{
		r:  r,
		tr: tr,
	}
}

func (u productUsecaseImpl) GetUserProducts(c context.Context, queryParams queryparams.QueryParams, userID int) (*entity.ProductsList, error) {
	if userID <= 0 {
		return u.GetGeneralProducts(c, queryParams)
	}

	isExists, err := u.r.IsAddressExistsByUserID(c, userID)
	if err != nil {
		return nil, err
	}

	if !isExists {
		return u.GetGeneralProducts(c, queryParams)
	}

	addressID, err := u.r.GetAddressByUserID(c, userID)
	if err != nil {
		return nil, err
	}

	location, err := u.r.GetLocationByAddressID(c, addressID)
	if err != nil {
		return nil, err
	}

	if queryParams.Limit <= 0 {
		queryParams.Limit = 10
	}

	totalProduct, products, err := u.r.GetUserProducts(c, queryParams, location)
	if err != nil {
		return nil, err
	}

	totalPage := int(math.Ceil(float64(totalProduct) / float64(queryParams.Limit)))
	if totalProduct <= 0 {
		totalPage = 1
	}
	
	pagination := entity.Pagination{
		CurrentPage:  queryParams.Page,
		TotalPage:    totalPage,
		TotalProduct: totalProduct,
	}
	productsList := entity.ProductsList{
		Pagination: pagination,
		Products:   products,
	}

	return &productsList, nil
}

func (u productUsecaseImpl) GetGeneralProducts(c context.Context, queryParams queryparams.QueryParams) (*entity.ProductsList, error) {
	var location = appconstant.DefaultLocation

	if queryParams.Limit <= 0 {
		queryParams.Limit = 10
	}

	totalProduct, products, err := u.r.GetUserProducts(c, queryParams, location)
	if err != nil {
		return nil, err
	}

	totalPage := int(math.Ceil(float64(totalProduct) / float64(queryParams.Limit)))
	if totalProduct <= 0 {
		totalPage = 1
	}

	pagination := entity.Pagination{
		CurrentPage:  queryParams.Page,
		TotalPage:    totalPage,
		TotalProduct: totalProduct,
	}
	productsList := entity.ProductsList{
		Pagination: pagination,
		Products:   products,
	}

	return &productsList, nil
}

func (u productUsecaseImpl) GetUserProductsHomePage(c context.Context, queryParams queryparams.QueryParams, userID int) (*entity.ProductsList, error) {
	if userID <= 0 {
		return u.GetGeneralProductsHomePage(c, queryParams)
	}

	isExists, err := u.r.IsAddressExistsByUserID(c, userID)
	if err != nil {
		return nil, err
	}

	if !isExists {
		return u.GetGeneralProductsHomePage(c, queryParams)
	}

	addressID, err := u.r.GetAddressByUserID(c, userID)
	if err != nil {

		return nil, err
	}

	location, err := u.r.GetLocationByAddressID(c, addressID)
	if err != nil {
		return nil, err
	}

	totalProduct, err := u.r.GetTotalProductHomePage(c, queryParams, location)
	if err != nil {
		return nil, err
	}

	if queryParams.Limit <= 0 {
		queryParams.Limit = 10
	}

	totalPage := int(math.Ceil(float64(totalProduct) / float64(queryParams.Limit)))
	if totalProduct <= 0 {
		totalPage = 1
	}

	queryParams.Page = apperror.CheckCurrentPage(queryParams.Page, totalPage, queryParams.Limit)

	products, err := u.r.GetUserProductsHomePage(c, queryParams, location)
	if err != nil {
		return nil, err
	}

	pagination := entity.Pagination{
		CurrentPage:  queryParams.Page,
		TotalPage:    totalPage,
		TotalProduct: totalProduct,
	}
	productsList := entity.ProductsList{
		Pagination: pagination,
		Products:   products,
	}

	return &productsList, nil
}

func (u productUsecaseImpl) GetGeneralProductsHomePage(c context.Context, queryParams queryparams.QueryParams) (*entity.ProductsList, error) {
	var location = appconstant.DefaultLocation

	totalProduct, err := u.r.GetTotalProductHomePage(c, queryParams, location)
	if err != nil {
		return nil, err
	}

	if queryParams.Limit <= 0 {
		queryParams.Limit = 10
	}

	totalPage := int(math.Ceil(float64(totalProduct) / float64(queryParams.Limit)))
	if totalProduct <= 0 {
		totalPage = 1
	}

	queryParams.Page = apperror.CheckCurrentPage(queryParams.Page, totalPage, queryParams.Limit)

	products, err := u.r.GetUserProductsHomePage(c, queryParams, location)
	if err != nil {
		return nil, err
	}

	pagination := entity.Pagination{
		CurrentPage:  queryParams.Page,
		TotalPage:    totalPage,
		TotalProduct: totalProduct,
	}
	productsList := entity.ProductsList{
		Pagination: pagination,
		Products:   products,
	}

	return &productsList, nil
}

func (u productUsecaseImpl) GetProductDetail(c context.Context, pharmacyProductID int) (*entity.ProductDetail, error) {
	isExists, err := u.r.IsPharmacyProductExistsByID(c, pharmacyProductID)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, apperror.NewErrStatusNotFound(appconstant.FieldErrGetProductDetail, apperror.ErrPharmacyProductNotExists, apperror.ErrPharmacyProductNotExists)
	}

	productDetail, err := u.r.GetProductDetail(c, pharmacyProductID)
	if err != nil {
		return nil, err
	}


	categories, err := u.r.GetProductCategories(c, productDetail.ID)
	if err != nil {
		return nil, err
	}
	productDetail.ProductCategories = categories

	return productDetail, nil
}

func (u productUsecaseImpl) GetMasterProducts(c context.Context, queryParams queryparams.QueryParams) (*entity.ProductsList, error) {
	queryparams.DefaultQueryParams(&queryParams)

	totalProduct, err := u.r.GetTotalMasterProduct(c, queryParams)
	if err != nil {
		return nil, err
	}

	if queryParams.Limit <= 0 {
		queryParams.Limit = 10
	}

	totalPage := int(math.Ceil(float64(totalProduct) / float64(queryParams.Limit)))
	if totalProduct <= 0 {
		totalPage = 1
	}

	queryParams.Page = apperror.CheckCurrentPage(queryParams.Page, totalPage, queryParams.Limit)

	products, err := u.r.GetMasterProducts(c, queryParams)
	if err != nil {
		return nil, err
	}

	pagination := entity.Pagination{
		CurrentPage:  queryParams.Page,
		TotalPage:    totalPage,
		TotalProduct: totalProduct,
	}

	productsList := entity.ProductsList{
		Pagination: pagination,
		Products:   products,
	}

	return &productsList, nil
}
