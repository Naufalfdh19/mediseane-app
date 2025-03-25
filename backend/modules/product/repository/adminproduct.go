package repository

import (
	"context"
	"fmt"
	"montelukast/modules/product/entity"
	queryparams "montelukast/modules/product/queryparams"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"
	"strings"
)

func (r ProductRepoImpl) GetProductsAdmin(c context.Context, queryParams queryparams.AdminQueryParams) ([]entity.ProductAdmin, error) {
	products := []entity.ProductAdmin{}

	query := `SELECT p.id, pc.name , pf.name, p.name, generic_name , manufacture, is_active,
				CASE 
					WHEN product_used IS NULL THEN 0
					WHEN product_used IS NOT NULL THEN product_used
				END AS product_used
				FROM products p 
				LEFT JOIN (
					SELECT product_id, count(*) AS product_used
					FROM pharmacy_products pp 
					GROUP BY product_id 
				) AS p1 ON p1.product_id = p.id 
				LEFT JOIN product_classifications pc ON pc.id = p.product_classification_id
				LEFT JOIN product_forms pf ON pf.id = p.product_form_id 
				WHERE p.deleted_at IS NULL AND pc.deleted_at IS NULL AND pf.deleted_at IS NULL`

	var params []any
	query += queryparams.AddAdminQueryParams(&params, queryParams)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	var product_used string
	for rows.Next() {
		var product entity.ProductAdmin
		err := rows.Scan(
			&product.ID,
			&product.ProductClassification,
			&product.ProductForm,
			&product.Name,
			&product.GenericName,
			&product.Manufacture,
			&product.IsActive,
			&product_used,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r ProductRepoImpl) GetTotalProductAdmin(c context.Context, queryParams queryparams.AdminQueryParams) (int, error) {
	queryParams.Limit = 0
	queryParams.Page = 0
	queryParams.SortBy = ""
	queryParams.Order = ""

	query := `SELECT COUNT(*)
				FROM products p 
				LEFT JOIN (
					SELECT product_id, count(*) AS product_used
					FROM pharmacy_products pp 
					GROUP BY product_id 
				) AS p1 ON p1.product_id = p.id 
				LEFT JOIN product_classifications pc ON pc.id = p.product_classification_id
				LEFT JOIN product_forms pf ON pf.id = p.product_form_id 
				WHERE p.deleted_at IS NULL AND pc.deleted_at IS NULL AND pf.deleted_at IS NULL`

	var params []any
	query += queryparams.AddAdminQueryParams(&params, queryParams)

	var totalProduct int
	err := r.db.QueryRow(query, params...).Scan(&totalProduct)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return totalProduct, nil
}

func (r ProductRepoImpl) UpdateProduct(c context.Context, product entity.Product) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE products
				SET product_classification_id = $2, product_form_id = $3,
					name = $4, generic_name = $5, manufacture = $6, description = $7,
					unit_in_pack = $8, weight = $9, height = $10, length = $11,
					width = $12, is_active = $13, updated_at = NOW()
				WHERE id = $1`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, product.ID,
			product.ProductClassificationID, product.ProductFormID,
			product.Name, product.GenericName, product.Manufacture, product.Description,
			product.UnitInPack, product.Weight, product.Height, product.Length, product.Width,
			product.IsActive,
		)
	} else {
		_, err = r.db.ExecContext(c, query, product.ID,
			product.ProductClassificationID, product.ProductFormID,
			product.Name, product.GenericName, product.Manufacture, product.Description,
			product.UnitInPack, product.Weight, product.Height, product.Length, product.Width,
			product.IsActive,
		)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r ProductRepoImpl) DeleteProduct(c context.Context, productID int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE products
               SET deleted_at = NOW()
               WHERE id = $1 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, productID)
	} else {
		_, err = r.db.ExecContext(c, query, productID)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r ProductRepoImpl) DeleteMultiCategories(c context.Context, productID int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE product_multi_categories
				SET deleted_at = NOW()
				WHERE product_id = $1 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, productID)
	} else {
		_, err = r.db.ExecContext(c, query, productID)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r ProductRepoImpl) DeletePharmacyDeletedProducts(c context.Context, productID int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE pharmacy_products
				SET deleted_at = NOW()
				WHERE product_id = $1 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, productID)
	} else {
		_, err = r.db.ExecContext(c, query, productID)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r ProductRepoImpl) AddMultipleCategories(c context.Context, product entity.Product) error {
	tx := transaction.ExtractTx(c)

	query := `INSERT INTO product_multi_categories (product_id, product_category_id)
			  VALUES `

	var params []any
	query += addCategories(&params, product)

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, params...)
	} else {
		_, err = r.db.ExecContext(c, query, params...)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func addCategories(params *[]any, product entity.Product) string {
	query := ``
	for paramIndex, category := range product.ProductCategoriesID {
		query += fmt.Sprintf("(%d, $%d), ", product.ID, paramIndex+1)
		*params = append(*params, category)
		paramIndex++
	}
	query = strings.TrimSuffix(query, ", ")
	return query
}

func (r ProductRepoImpl) UpdateProductPhoto(c context.Context, url string, productID int) error {
	query := `UPDATE product
				SET image[1] = $2, updated_at = NOW()
				WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(query, productID, url)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r ProductRepoImpl) GetTotalProductCategories(c context.Context, categories []int) (int, error) {
	var params []any
	querIndex := 1

	query := `SELECT COUNT(*)
				FROM product_categories pc 
				WHERE deleted_at IS NULL AND`
	query += queryparams.AdminGetCategoriesParams(&params, categories, &querIndex)
	var totalProductCategories int
	err := r.db.QueryRow(query, params...).Scan(&totalProductCategories)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return totalProductCategories, nil
}
