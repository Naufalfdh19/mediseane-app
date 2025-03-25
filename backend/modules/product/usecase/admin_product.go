package usecase

import (
	"context"
	"math"
	"montelukast/modules/product/entity"
	queryparams "montelukast/modules/product/queryparams"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/imageuploader"

	"github.com/go-playground/validator/v10"
)

func (u productUsecaseImpl) AddProduct(c context.Context, product entity.Product) error {
	err := u.tr.WithinTransaction(c, func(TxCtx context.Context) error {
		err := CheckProductClass(product)
		if err != nil {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdateProduct, err, err)
		}

		totalProductCategory, err := u.r.GetTotalProductCategories(TxCtx, product.ProductCategoriesID)
		if err != nil {
			return err
		}
		if totalProductCategory != len(product.ProductCategoriesID) {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdateProduct, apperror.ErrCategoryNotExists, apperror.ErrCategoryNotExists)
		}

		err = checkProductCategory(product)
		if err != nil {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddProduct, err, err)
		}

		isProductExists, err := u.r.IsProductExists(TxCtx, product)
		if err != nil {
			return err
		}
		if isProductExists {
			err := apperror.NewErrStatusBadRequest(appconstant.FieldErrAddProduct, apperror.ErrProductAlreadyExists, apperror.ErrProductAlreadyExists)
			return err
		}

		err = u.r.AddProduct(TxCtx, &product)
		if err != nil {
			return err
		}

		err = u.r.AddMultipleCategories(TxCtx, product)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u productUsecaseImpl) DeleteProduct(c context.Context, productID int) error {
	err := u.tr.WithinTransaction(c, func(txCtx context.Context) error {
		isProductExists, err := u.r.IsProductExistsByID(c, productID)
		if err != nil {
			return err
		}
		if !isProductExists {
			return apperror.NewErrStatusNotFound(appconstant.FieldErrDeleteProduct, apperror.ErrProductNotExists, apperror.ErrProductNotExists)
		}

		err = u.r.DeleteProduct(txCtx, productID)
		if err != nil {
			return err
		}

		err = u.r.DeleteMultiCategories(txCtx, productID)
		if err != nil {
			return err
		}

		err = u.r.DeletePharmacyDeletedProducts(txCtx, productID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u productUsecaseImpl) UpdateProduct(c context.Context, product entity.Product) error {
	err := u.tr.WithinTransaction(c, func(txCtx context.Context) error {
		err := CheckProductClass(product)
		if err != nil {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdateProduct, err, err)
		}

		totalProductCategory, err := u.r.GetTotalProductCategories(txCtx, product.ProductCategoriesID)
		if err != nil {
			return err
		}
		if totalProductCategory != len(product.ProductCategoriesID) {
			return apperror.NewErrStatusNotFound(appconstant.FieldErrUpdateProduct, apperror.ErrCategoryNotExists, apperror.ErrCategoryNotExists)
		}

		err = checkProductCategory(product)
		if err != nil {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdateProduct, err, err)
		}

		isProductExists, err := u.r.IsProductExistsByID(txCtx, product.ID)
		if err != nil {
			return err
		}
		if !isProductExists {
			return apperror.NewErrStatusNotFound(appconstant.FieldErrUpdateProduct, apperror.ErrProductNotExists, apperror.ErrProductNotExists)
		}

		err = u.r.UpdateProduct(txCtx, product)
		if err != nil {
			return err
		}

		err = u.r.DeleteMultiCategories(txCtx, product.ID)
		if err != nil {
			return err
		}

		err = u.r.AddMultipleCategories(txCtx, product)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func CheckProductClass(product entity.Product) error {
	if product.ProductClassificationID > 4 || product.ProductClassificationID < 1 {
		return apperror.ErrProductClassNotExists
	}

	if product.ProductClassificationID == 4 {
		return nil
	}
	if product.ProductFormID == nil {
		return apperror.ErrProductFormMandatory
	}
	if product.UnitInPack == nil {
		return apperror.ErrUnitInPackMandatory
	}
	return nil
}

func checkProductCategory(product entity.Product) error {
	categories := map[int]bool{}
	for _, categoryID := range product.ProductCategoriesID {
		if _, exists := categories[categoryID]; exists {
			return apperror.ErrDuplicateCategory
		}
		categories[categoryID] = true
	}

	return nil
}

func (u productUsecaseImpl) GetProductsAdmin(c context.Context, queryParams queryparams.AdminQueryParams) (*entity.ProductListAdmin, error) {
	totalProduct, err := u.r.GetTotalProductAdmin(c, queryParams)
	if err != nil {
		return nil, err
	}

	queryParams = queryparams.AdminDefaultQuery(queryParams)

	totalPage := int(math.Ceil(float64(totalProduct) / float64(queryParams.Limit)))
	if totalProduct <= 0 {
		totalPage = 1
	}

	queryParams.Page = apperror.CheckCurrentPage(queryParams.Page, totalPage, queryParams.Limit)

	products, err := u.r.GetProductsAdmin(c, queryParams)
	if err != nil {
		return nil, err
	}

	pagination := entity.Pagination{
		CurrentPage:  queryParams.Page,
		TotalPage:    totalPage,
		TotalProduct: totalProduct,
	}

	productsList := entity.ProductListAdmin{
		Pagination: pagination,
		Products:   products,
	}

	return &productsList, nil
}

func (u productUsecaseImpl) UpdateProductPhoto(c context.Context, file entity.File, productID int) (string, error) {
	isPharmacistExists, err := u.r.IsProductExistsByID(c, productID)
	if err != nil {
		return "", err
	}
	if !isPharmacistExists {
		return "", apperror.NewErrStatusNotFound(appconstant.FieldUpdateProductPhoto, apperror.ErrProductNotExists, apperror.ErrProductNotExists)
	}

	validate := validator.New()
	err = validate.Struct(file)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrUploadImage, apperror.ErrInternalServer, err)
	}
	uploadUrl, err := imageuploader.ImageUploadHelper(file.File)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrUploadImage, apperror.ErrInternalServer, err)
	}
	err = u.r.UpdateProductPhoto(c, uploadUrl, productID)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}
