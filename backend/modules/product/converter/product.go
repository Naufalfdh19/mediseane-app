package converter

import (
	queryparams "montelukast/modules/product/queryparams"
	"montelukast/modules/product/dto"
	"montelukast/modules/product/entity"
)

type GetUserProductsConverter struct{}

func (c GetUserProductsConverter) ToDto(product entity.ProductDetail) dto.GetProductsResponse {
	return dto.GetProductsResponse{
		ID:                product.ID,
		PharmacyProductID: product.PharmacyProductID,
		Name:              product.Name,
		Manufacture:       product.Manufacture,
		Image:             product.Image,
		PharmacyName:      product.PharmacyName,
		Price:             product.Price,
	}
}

type GetListOfProductsConverter struct{}

func (c GetListOfProductsConverter) ToDto(product entity.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ProductCategoriesID:     product.ProductCategoriesID,
		ProductClassificationID: product.ProductClassificationID,
		ProductFormID:           product.ProductFormID,
		Name:                    product.Name,
		GenericName:             product.GenericName,
		Manufacture:             product.Manufacture,
		Description:             product.Description,
		Image:                   product.Image,
		UnitInPack:              product.UnitInPack,
		Weight:                  product.Weight,
		Height:                  product.Height,
		Length:                  product.Length,
		Width:                   product.Width,
	}
}

type AddProductsConverter struct{}

func (c AddProductsConverter) ToEntity(product dto.AddProductRequest) entity.Product {
	return entity.Product{
		ProductCategoriesID:     product.ProductCategoriesID,
		ProductClassificationID: product.ProductClassificationID,
		ProductFormID:           product.ProductFormID,
		Name:                    product.Name,
		GenericName:             product.GenericName,
		Manufacture:             product.Manufacture,
		Description:             product.Description,
		UnitInPack:              product.UnitInPack,
		Weight:                  product.Weight,
		Height:                  product.Height,
		Length:                  product.Length,
		Width:                   product.Width,
		IsActive:                product.IsActive,
	}
}

type PaginationConverter struct{}

func (c PaginationConverter) ToDto(pagination entity.Pagination) dto.Pagination {
	return dto.Pagination{
		CurrentPage:  pagination.CurrentPage,
		TotalProduct: pagination.TotalProduct,
		TotalPage:    pagination.TotalPage,
	}
}

type GetProductDetailConverter struct{}

func (c GetProductDetailConverter) ToEntity(product dto.GetProductDetailRequest) entity.Product {
	return entity.Product{
		PharmacyProductID: product.PharmacistProductID,
	}
}

type ProductDetailConverter struct{}

func (c ProductDetailConverter) ToDto(productDetail entity.ProductDetail) dto.GetProductDetailResponse {
	return dto.GetProductDetailResponse{
		ID:                productDetail.ID,
		PharmacyProductID: productDetail.PharmacyProductID,
		ProductCategories: productDetail.ProductCategories,
		Name:              productDetail.Name,
		GenericName:       productDetail.GenericName,
		Manufacture:       productDetail.Manufacture,
		Description:       productDetail.Description,
		Image:             productDetail.Image,
		UnitInPack:        *productDetail.UnitInPack,
		Stock:             productDetail.Stock,
		Price:             productDetail.Price,
		PharmacyAddress:   productDetail.PharmacyAddress,
		PharmacysName:     productDetail.PharmacyName,
	}
}

type AdminQueryParamsConverter struct{}

func (c AdminQueryParamsConverter) ToEntity(queryParams queryparams.AdminQueryParamsDto) queryparams.AdminQueryParams {
	return queryparams.AdminQueryParams{
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

type QueryParamsConverter struct{}

func (c QueryParamsConverter) ToEntity(queryParams queryparams.QueryParamsDto) queryparams.QueryParams {
	return queryparams.QueryParams{
		Name:       queryParams.Name,
		CategoryID: queryParams.CategoryID,
		SortBy:     queryParams.SortBy,
		Order:      queryParams.Order,
		Limit:      queryParams.Limit,
		Page:       queryParams.Page,
	}
}

type GetProductsAdminConverter struct{}

func (c GetProductsAdminConverter) ToDto(product entity.ProductAdmin) dto.GetProductResponse {
	return dto.GetProductResponse{
		ID:                    product.ID,
		ProductClassification: product.ProductClassification,
		ProductForm:           product.ProductForm,
		Name:                  product.Name,
		GenericName:           product.GenericName,
		Manufacture:           product.Manufacture,
		IsActive:              product.IsActive,
	}
}

type FileConverter struct{}

func (c FileConverter) ToEntity(fileDTO dto.FileRequest) entity.File {
	return entity.File{
		File: fileDTO.File,
	}
}

type GetMasterProductConverter struct {}

func (c GetMasterProductConverter) ToDto(product entity.ProductDetail) dto.GetMasterProductResponse {
	return dto.GetMasterProductResponse{
		ID: product.ID,
		Name: product.Name,
	}
}
