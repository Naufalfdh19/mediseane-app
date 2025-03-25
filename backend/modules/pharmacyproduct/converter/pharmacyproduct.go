package converter

import (
	"montelukast/modules/pharmacyproduct/dto"
	"montelukast/modules/pharmacyproduct/entity"
	"montelukast/modules/pharmacyproduct/queryparams"
	productEntity "montelukast/modules/product/entity"
)

type AddPharmacyProductConverter struct{}

func (c AddPharmacyProductConverter) ToEntity(pharmacyProduct dto.AddPharmacyProductRequest) entity.PharmacyProduct {
	return entity.PharmacyProduct{
		ProductID: pharmacyProduct.ProductID,
		Stock:     pharmacyProduct.Stock,
		Price:     pharmacyProduct.Price,
		IsActive:  pharmacyProduct.IsActive,
	}
}

type UpdatePharmacyProductConverter struct{}

func (c UpdatePharmacyProductConverter) ToEntity(pharmacistProduct dto.UpdatePharmacyProductRequest) entity.PharmacyProduct {
	return entity.PharmacyProduct{
		Stock:    pharmacistProduct.Stock,
		IsActive: pharmacistProduct.IsActive,
	}
}

type GetPharmacyProductConverter struct{}

func (c GetPharmacyProductConverter) ToDto(pharmacistProduct productEntity.ProductDetail) dto.GetPharmacyProductResponse {
	return dto.GetPharmacyProductResponse{
		ID:                    pharmacistProduct.PharmacyProductID,
		Name:                  pharmacistProduct.Name,
		GenericName:           pharmacistProduct.GenericName,
		Manufacturer:          pharmacistProduct.Manufacture,
		ProductClassification: pharmacistProduct.ProductClassification,
		ProductForm:           pharmacistProduct.ProductForm,
		Stock:                 pharmacistProduct.Stock,
		Price:                 pharmacistProduct.Price.String(),
		IsActive:              pharmacistProduct.IsActive,
		CreatedAt:             pharmacistProduct.CreatedAt,
	}
}

type QueryParamsConverter struct{}

func (c QueryParamsConverter) ToEntity(queryParams queryparams.QueryParamsDto) queryparams.QueryParams {
	return queryparams.QueryParams{
		Name:                  queryParams.Name,
		GenericName:           queryParams.GenericName,
		Manufacture:           queryParams.Manufacture,
		ProductClassification: queryParams.ProductClassification,
		ProductForm:           queryParams.ProductForm,
		IsActive:              queryParams.IsActive,
		SortBy:                queryParams.SortBy,
		Order:                 queryParams.Order,
		Limit:                 queryParams.Limit,
		Page:                  queryParams.Page,
	}
}
