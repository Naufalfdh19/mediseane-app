package usecase

import (
	"context"
	"math"
	pharmacistRepo "montelukast/modules/pharmacist/repository"
	pharmacyRepo "montelukast/modules/pharmacy/repository"
	"montelukast/modules/pharmacyproduct/entity"
	"montelukast/modules/pharmacyproduct/queryparams"
	"montelukast/modules/pharmacyproduct/repository"
	productEntity "montelukast/modules/product/entity"
	productRepo "montelukast/modules/product/repository"
	appconstant "montelukast/pkg/constant"
	"montelukast/pkg/dateconverter"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
)

type PharmacyProductUsecase interface {
	AddPharmacyProduct(c context.Context, pharmacyProduct entity.PharmacyProduct, pharmacistID int) error
	UpdatePharmacyProduct(c context.Context, pharmacyProduct entity.PharmacyProduct, pharmacistID int) error
	DeletePharmacyProduct(c context.Context, pharmacyProduct entity.PharmacyProduct, pharmacistID int) error
	GetPharmacyProductDetail(c context.Context, pharmacyProductID int, pharmacistID int) (*productEntity.ProductDetail, error)
	GetPharmacyProducts(c context.Context, queryParams queryparams.QueryParams, pharmacistID int) (*productEntity.ProductsList, error)
}

type pharmacyProductUsecaseImpl struct {
	r    repository.PharmacyProductRepo
	tr   transaction.TransactorRepoImpl
	phr  pharmacyRepo.PharmacyRepository
	phsr pharmacistRepo.PharmacistRepo
	pr   productRepo.ProductRepo
}

func NewPharmacyProductUsecase(r repository.PharmacyProductRepo, tr transaction.TransactorRepoImpl, phr pharmacyRepo.PharmacyRepository, pr productRepo.ProductRepo, phsr pharmacistRepo.PharmacistRepo) pharmacyProductUsecaseImpl {
	return pharmacyProductUsecaseImpl{
		r:    r,
		tr:   tr,
		phr:  phr,
		phsr: phsr,
		pr:   pr,
	}
}

func (u pharmacyProductUsecaseImpl) AddPharmacyProduct(c context.Context, pharmacyProduct entity.PharmacyProduct, pharmacistID int) error {
	if pharmacyProduct.Stock < 0 || decimal.NewFromInt(0).GreaterThan(pharmacyProduct.Price) {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPharmacyProduct, apperror.ErrPriceOrStockLessThanZero, apperror.ErrPriceOrStockLessThanZero)
	}

	isExists, err := u.phsr.IsPharmacistExistsByID(c, pharmacistID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPharmacyProduct, apperror.ErrPharmacistNotExists, apperror.ErrPharmacistNotExists)
	}

	pharmacyID, err := u.r.GetPharmacyIDbyPharmacistID(c, pharmacistID)
	if err != nil {
		return err
	}
	if pharmacyID == 0 {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPharmacyProduct, apperror.ErrPharmacistNotHasPharmacy, apperror.ErrPharmacistNotHasPharmacy)
	}
	pharmacyProduct.PharmacyID = pharmacyID

	isExists, err = u.r.IsPharmacyProductExists(c, pharmacyID, pharmacyProduct.ProductID)
	if err != nil {
		return err
	}
	if isExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPharmacyProduct, apperror.ErrProductAlreadyExists, apperror.ErrProductAlreadyExists)
	}

	isExists, err = u.phr.IsPharmacyExists(c, pharmacyID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrAddPharmacyProduct, apperror.ErrPharmacyNotExists, apperror.ErrPharmacyNotExists)
	}

	isExists, err = u.r.IsProductExistsByID(c, pharmacyProduct.ProductID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrAddPharmacyProduct, apperror.ErrProductNotExists, apperror.ErrProductNotExists)
	}

	err = u.r.AddPharmacyProduct(c, pharmacyProduct)
	if err != nil {
		return err
	}

	return nil
}

func (u pharmacyProductUsecaseImpl) UpdatePharmacyProduct(c context.Context, pharmacyProduct entity.PharmacyProduct, pharmacistID int) error {
	isExists, err := u.phsr.IsPharmacistExistsByID(c, pharmacistID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPharmacyProduct, apperror.ErrPharmacistNotExists, apperror.ErrPharmacistNotExists)
	}

	isExists, err = u.r.IsPharmacyProductExistsByID(c, pharmacyProduct.ID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrUpdatePharmacyProduct, apperror.ErrPharmacyProductNotExists, apperror.ErrPharmacyProductNotExists)
	}

	err = u.CheckPharmacistAuthorization(c, pharmacistID, pharmacyProduct.ID)
	if err != nil {
		return nil
	}

	currStock, err := u.r.GetStockByID(c, pharmacyProduct.ID)
	if err != nil {
		return err
	}

	err = u.CheckStockUpdatedDate(c, currStock, pharmacyProduct.Stock)
	if err != nil {
		return err
	}

	err = u.r.UpdatePharmacyProduct(c, pharmacyProduct)
	if err != nil {
		return err
	}

	err = u.r.SetStockUpdatedDate(c)
	if err != nil {
		return err
	}

	return nil
}

func (u pharmacyProductUsecaseImpl) DeletePharmacyProduct(c context.Context, pharmacyProduct entity.PharmacyProduct, pharmacistID int) error {
	isExists, err := u.phsr.IsPharmacistExistsByID(c, pharmacistID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddPharmacyProduct, apperror.ErrPharmacistNotExists, apperror.ErrPharmacistNotExists)
	}

	isExists, err = u.r.IsPharmacyProductExistsByID(c, pharmacyProduct.ID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrUpdatePharmacyProduct, apperror.ErrPharmacyProductNotExists, apperror.ErrPharmacyProductNotExists)
	}

	err = u.CheckPharmacistAuthorization(c, pharmacistID, pharmacyProduct.ID)
	if err != nil {
		return nil
	}

	err = u.r.DeletePharmacyProduct(c, pharmacyProduct.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u pharmacyProductUsecaseImpl) GetPharmacyProductDetail(c context.Context, pharmacyProductID int, pharmacistID int) (*productEntity.ProductDetail, error) {
	isExists, err := u.phsr.IsPharmacistExistsByID(c, pharmacistID)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, apperror.NewErrStatusBadRequest(appconstant.FieldErrGetProduct, apperror.ErrPharmacistNotExists, apperror.ErrPharmacistNotExists)
	}

	pharmacyID, err := u.r.GetPharmacyIDbyPharmacistID(c, pharmacistID)
	if err != nil {
		return nil, err
	}
	if pharmacyID == 0 {
		return nil, apperror.NewErrStatusBadRequest(appconstant.FieldErrCheckAuthorization, apperror.ErrPharmacistNotHasPharmacy, apperror.ErrPharmacistNotHasPharmacy)
	}

	isExists, err = u.r.IsPharmacyProductExistsByIDAndPharmacy(c, pharmacyProductID, pharmacyID)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, apperror.NewErrStatusNotFound(appconstant.FieldErrGetProduct, apperror.ErrPharmacyProductNotExists, apperror.ErrPharmacyProductNotExists)
	}

	product, err := u.r.GetPharmacyProduct(c, pharmacyProductID, pharmacyID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u pharmacyProductUsecaseImpl) GetPharmacyProducts(c context.Context, queryParams queryparams.QueryParams, pharmacistID int) (*productEntity.ProductsList, error) {
	queryparams.DefaultQueryParams(&queryParams)

	isExists, err := u.phsr.IsPharmacistExistsByID(c, pharmacistID)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, apperror.NewErrStatusBadRequest(appconstant.FieldErrGetPharmacyProducts, apperror.ErrPharmacistNotExists, apperror.ErrPharmacistNotExists)
	}

	pharmacyID, err := u.r.GetPharmacyIDbyPharmacistID(c, pharmacistID)
	if err != nil {
		return nil, err
	}
	if pharmacyID == 0 {
		return nil, apperror.NewErrStatusBadRequest(appconstant.FieldErrCheckAuthorization, apperror.ErrPharmacistNotHasPharmacy, apperror.ErrPharmacistNotHasPharmacy)
	}

	totalProduct, err := u.r.GetTotalProduct(c, queryParams, pharmacyID)
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

	products, err := u.r.GetPharmacyProducts(c, queryParams, pharmacyID)
	if err != nil {
		return nil, err
	}

	pagination := productEntity.Pagination{
		CurrentPage:  queryParams.Page,
		TotalPage:    totalPage,
		TotalProduct: totalProduct,
	}

	productsList := productEntity.ProductsList{
		Pagination: pagination,
		Products:   products,
	}

	return &productsList, nil
}

func (u pharmacyProductUsecaseImpl) CheckStockUpdatedDate(c context.Context, currStock, incomingStock int) error {
	if currStock != incomingStock {
		stockUpdatedDate, err := u.r.GetStockUpdatedDate(c)
		if err != nil && err.Error() != redis.Nil.Error() {
			return err
		}
		if err != redis.Nil {
			currentDate := dateconverter.GetCurrentDate()
			if stockUpdatedDate == currentDate {
				return apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdatePharmacyProduct, apperror.ErrStockAlreadyUpdated, apperror.ErrStockAlreadyUpdated)
			}
		}
	}

	return nil
}

func (u pharmacyProductUsecaseImpl) CheckPharmacistAuthorization(c context.Context, pharmacistID, pharmacyProductID int) error {
	pharmacyID1, err := u.r.GetPharmacyIDbyPharmacistID(c, pharmacistID)
	if err != nil {
		return err
	}
	if pharmacyID1 == 0 {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrCheckAuthorization, apperror.ErrPharmacistNotHasPharmacy, apperror.ErrPharmacistNotHasPharmacy)
	}

	pharmacyID2, err := u.r.GetPharmacyIDbyPharmacyProductID(c, pharmacyProductID)
	if err != nil {
		return err
	}

	if pharmacyID1 != pharmacyID2 {
		return apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrPharmacistUnauthorized, apperror.ErrPharmacistUnauthorized)
	}

	return nil
}
