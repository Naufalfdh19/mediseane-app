package repository

import (
	"context"
	"database/sql"

	"montelukast/modules/product/entity"
	queryparams "montelukast/modules/product/queryparams"

	appconstant "montelukast/pkg/constant"

	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"
)

type ProductRepo interface {
	GetUserProducts(c context.Context, queryParams queryparams.QueryParams, location string) (int, []entity.ProductDetail, error)
	IsAddressExistsByUserID(c context.Context, userID int) (bool, error)
	GetAddressByUserID(c context.Context, userID int) (int, error)
	GetProductDetail(c context.Context, pharcistsProductID int) (*entity.ProductDetail, error)
	GetProductsAdmin(c context.Context, queryParams queryparams.AdminQueryParams) ([]entity.ProductAdmin, error)
	GetTotalProductAdmin(c context.Context, queryParams queryparams.AdminQueryParams) (int, error)
	IsProductExistsByID(c context.Context, productID int) (bool, error)
	UpdateProduct(c context.Context, product entity.Product) error
	DeleteProduct(c context.Context, productID int) error
	AddMultipleCategories(c context.Context, product entity.Product) error
	DeleteMultiCategories(c context.Context, productID int) error
	GetProductCategories(c context.Context, productID int) ([]string, error)
	GetLocationByAddressID(c context.Context, addressID int) (string, error)
	GetTotalProductHomePage(c context.Context, queryParams queryparams.QueryParams, location string) (int, error)
	GetUserProductsHomePage(c context.Context, queryParams queryparams.QueryParams, location string) ([]entity.ProductDetail, error)
	IsPharmacyProductExistsByID(c context.Context, id int) (bool, error)
	GetMasterProducts(c context.Context, queryParams queryparams.QueryParams) ([]entity.ProductDetail, error)
	GetTotalMasterProduct(c context.Context, queryParams queryparams.QueryParams) (int, error)
	AddProduct(c context.Context, product *entity.Product) error
	IsProductExists(c context.Context, product entity.Product) (bool, error)
	DeletePharmacyDeletedProducts(c context.Context, productID int) error
	UpdateProductPhoto(c context.Context, url string, productID int) error
	GetCategoryBoundary(c context.Context) (*entity.CategoryBoundary, error)
	GetTotalProductCategories(c context.Context, categories []int) (int, error)
}

type ProductRepoImpl struct {
	db *sql.DB
}

func NewProductRepo(dbConn *sql.DB) ProductRepoImpl {
	return ProductRepoImpl{
		db: dbConn,
	}
}

func (r ProductRepoImpl) GetUserProducts(c context.Context, queryParams queryparams.QueryParams, location string) (int, []entity.ProductDetail, error) {
	query := `SELECT 
							p.id,
							pfs.pharmacy_product_id,
							p.image,
							p.name,
							p.manufacture,
							pfs.pharmacy_name,
							pfs.product_price ,
							pfs.total_score AS total_score,
							count(p.id) OVER ()
						FROM (
							SELECT DISTINCT ON (product_id)
								product_id,
								pharmacy_product_id,
								partner_id,
								pharmacy_name,
								product_price,
								price_score + ((1 - (distance - min(distance) OVER ()) / (max(distance) OVER () - min(distance) OVER ())) * 0.3) AS total_score
							FROM (
								SELECT
									pp.product_id AS product_id,
									pp.id AS pharmacy_product_id,
									p2.partner_id AS partner_id,
									p2.name AS pharmacy_name,
									pp.price AS product_price,
									(1 - (pp.price - min(pp.price) OVER ()) / (max(pp.price) OVER () - min(pp.price) OVER ())) * 0.7 AS price_score,
									ST_Distance(p2.location, $1::geometry) / 1000 AS distance
								FROM pharmacies p2 
								JOIN pharmacy_products pp 
									ON p2.id = pp.pharmacy_id
									AND pp.deleted_at IS NULL
									AND ST_DWithin(p2.location, $1::geometry, 25000)
									AND p2.deleted_at IS NULL
									AND p2.is_active = TRUE`

	var params []any
	params = append(params, location)

	querIndex := 2
	query += queryparams.AddFilterByCategoryID(&params, queryParams, &querIndex)

	query += `)
						ORDER BY product_id, total_score DESC) pfs
					JOIN partners ps
						ON ps.id = pfs.partner_id
						AND ps.active_days ILIKE '%' || TO_CHAR(CURRENT_DATE, 'FMDay') || '%'
						AND CURRENT_TIME BETWEEN ps.start_hour AND ps.end_hour
						AND ps.is_active = TRUE
						AND ps.deleted_at IS NULL
					JOIN products p
						ON pfs.product_id = p.id`

	query += queryparams.AddConditionQuery(&params, queryParams, &querIndex)
	query += ` ORDER BY total_score DESC`
	query += queryparams.AddPaginationQuery(&params, queryParams, &querIndex)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return -1, nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	products := []entity.ProductDetail{}
	var count int

	for rows.Next() {
		var product entity.ProductDetail
		var score float64
		err := rows.Scan(
			&product.ID,
			&product.PharmacyProductID,
			&product.Image,
			&product.Name,
			&product.Manufacture,
			&product.PharmacyName,
			&product.Price,
			&score,
			&count,
		)
		if err != nil {
			return -1, nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		products = append(products, product)
	}
	return count, products, nil
}

func (r ProductRepoImpl) GetTotalProductHomePage(c context.Context, queryParams queryparams.QueryParams, location string) (int, error) {
	query := `WITH most_bought AS (
							SELECT 
									p2.id AS product_id,
									p.partner_id AS partner_id
							FROM order_product_details opd
							JOIN pharmacy_products pp 
									ON opd.pharmacy_product_id = pp.id 
									AND opd.created_at >= (NOW()::date)
									AND opd.deleted_at IS NULL
									AND pp.deleted_at IS NULL
							JOIN pharmacies p 
									ON ST_DWithin(p.location, $1::geometry, 25000)
									AND pp.pharmacy_id = p.id 
									AND p.is_active = true
									AND p.deleted_at IS NULL
							JOIN products p2 
									ON pp.product_id = p2.id
									AND p2.deleted_at IS NULL)
						SELECT 
							count(DISTINCT (product_id))
						FROM most_bought mb
						JOIN partners ps
							ON mb.partner_id = ps.id
							AND ps.is_active = TRUE
							AND ps.deleted_at IS NULL
						WHERE 
							ps.active_days ILIKE '%' || TO_CHAR(CURRENT_DATE, 'FMDay') || '%'
						AND
							CURRENT_TIME BETWEEN ps.start_hour AND ps.end_hour;`

	var totalProduct int
	err := r.db.QueryRow(query, location).Scan(&totalProduct)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return totalProduct, nil
}

func (r ProductRepoImpl) GetUserProductsHomePage(c context.Context, queryParams queryparams.QueryParams, location string) ([]entity.ProductDetail, error) {
	products := []entity.ProductDetail{}

	query := `WITH most_bought AS (
							SELECT 
									p2.id AS product_id,
									pp.id AS pharmacy_product_id,
									p2.image AS	product_image,
									p2.name AS	product_name,
									p2.manufacture AS product_manufacture,
									p.name AS pharmacy_name,
									pp.price AS product_price,
									p.partner_id AS partner_id,
									(1 - (pp.price - min(pp.price) OVER ()) / (max(pp.price) OVER () - min(pp.price) OVER ())) * 0.7 AS price_score,
									ST_Distance(p.location, $1::geometry) / 1000 AS distance
							FROM order_product_details opd
							JOIN pharmacy_products pp 
									ON opd.pharmacy_product_id = pp.id 
									AND opd.created_at >= (NOW()::date)
									AND opd.deleted_at IS NULL
									AND pp.deleted_at IS NULL
							JOIN pharmacies p 
									ON ST_DWithin(p.location, $1::geometry, 25000)
									AND pp.pharmacy_id = p.id 
									AND p.is_active = true
									AND p.deleted_at IS NULL
							JOIN products p2 
									ON pp.product_id = p2.id
									AND p2.deleted_at IS NULL)
						SELECT * FROM (
							SELECT DISTINCT ON (product_id)
								product_id,
								pharmacy_product_id,
								product_image,
								product_name,
								product_manufacture,
								pharmacy_name,
								product_price,
								price_score + ((1 - (distance - min(distance) OVER ()) / (max(distance) OVER () - min(distance) OVER ())) * 0.3) AS total_score
							FROM most_bought mb
							JOIN partners ps
								ON mb.partner_id = ps.id
								AND ps.is_active = TRUE
								AND ps.deleted_at IS NULL
							WHERE 
								ps.active_days ILIKE '%' || TO_CHAR(CURRENT_DATE, 'FMDay') || '%'
							AND
								CURRENT_TIME BETWEEN ps.start_hour AND ps.end_hour
							ORDER BY product_id, total_score DESC)
						ORDER BY total_score DESC`

	var params []any
	querIndex := 2

	params = append(params, location)
	query += queryparams.AddPaginationQuery(&params, queryParams, &querIndex)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	var score float64
	for rows.Next() {
		var product entity.ProductDetail
		err := rows.Scan(
			&product.ID,
			&product.PharmacyProductID,
			&product.Image,
			&product.Name,
			&product.Manufacture,
			&product.PharmacyName,
			&product.Price,
			&score,
		)
		if err != nil {

			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		products = append(products, product)
	}
	return products, nil
}

func (r ProductRepoImpl) IsAddressExistsByUserID(c context.Context, userID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM user_addresses WHERE user_id = $1 AND is_active = true AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, userID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r ProductRepoImpl) GetAddressByUserID(c context.Context, userID int) (int, error) {
	query := `SELECT id
				FROM user_addresses
				WHERE user_id = $1 AND is_active = true AND deleted_at IS NULL`

	var addressID int
	err := r.db.QueryRow(query, userID).Scan(&addressID)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return addressID, nil
}

func (r ProductRepoImpl) GetProductCategories(c context.Context, productID int) ([]string, error) {
	categories := []string{}

	query := `select pc.name
				from products p 
				join product_multi_categories pmc on pmc.product_id = p.id 
				join product_categories pc on pc.id = pmc.product_category_id 
				where p.id = $1
				`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		err := rows.Scan(
			&category,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r ProductRepoImpl) GetProductDetail(c context.Context, pharcistsProductID int) (*entity.ProductDetail, error) {
	query := `select p.id, pp.id, p.name, p.image[1], p.generic_name, p.manufacture, p.description, p.unit_in_pack, ph.name, ph.address, pp.stock, pp.price 
				from pharmacy_products pp 
				join products p on p.id = pp.product_id and p.deleted_at is null
				join pharmacies ph on ph.id = pp.pharmacy_id and ph.deleted_at is null
				where pp.id = $1 and p.deleted_at is null`

	var productDetail entity.ProductDetail
	err := r.db.QueryRow(query, pharcistsProductID).Scan(
		&productDetail.ID,
		&productDetail.PharmacyProductID,
		&productDetail.Name,
		&productDetail.Image,
		&productDetail.GenericName,
		&productDetail.Manufacture,
		&productDetail.Description,
		&productDetail.UnitInPack,
		&productDetail.PharmacyName,
		&productDetail.PharmacyAddress,
		&productDetail.Stock,
		&productDetail.Price,
	)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return &productDetail, nil
}

func (r ProductRepoImpl) AddProduct(c context.Context, product *entity.Product) error {
	tx := transaction.ExtractTx(c)

	query := `INSERT INTO products (product_classification_id, product_form_id, name, generic_name, manufacture, description, image, unit_in_pack, weight, height, length, width, is_active)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			  RETURNING id`

	var err error

	if tx != nil {
		err = tx.QueryRowContext(c, query,
			product.ProductClassificationID, product.ProductFormID,
			product.Name, product.GenericName, product.Manufacture, product.Description,
			[]string{product.Image}, product.UnitInPack, product.Weight, product.Height, product.Length, product.Width,
			product.IsActive,
		).Scan(&product.ID)
	} else {
		err = r.db.QueryRowContext(c, query,
			product.ProductClassificationID, product.ProductFormID,
			product.Name, product.GenericName, product.Manufacture, product.Description,
			[]string{product.Image}, product.UnitInPack, product.Weight, product.Height, product.Length, product.Width,
			product.IsActive,
		).Scan(&product.ID)
	}

	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r ProductRepoImpl) IsProductExists(c context.Context, product entity.Product) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM products WHERE name = $1 AND generic_name = $2 AND manufacture = $3 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, product.Name, product.GenericName, product.Manufacture).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r ProductRepoImpl) IsProductExistsByID(c context.Context, productID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM products WHERE id = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, productID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r ProductRepoImpl) GetLocationByAddressID(c context.Context, addressID int) (string, error) {
	query := `SELECT location
				FROM user_addresses
				WHERE id = $1 AND is_active = true AND deleted_at IS NULL`

	var location string
	err := r.db.QueryRow(query, addressID).Scan(&location)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return location, nil
}

func (r ProductRepoImpl) IsPharmacyProductExistsByID(c context.Context, id int) (bool, error) {
	query := `SELECT EXISTS (
							SELECT 1 
							FROM pharmacy_products pp 
							JOIN products p 
							ON pp.product_id = p.id 
							WHERE pp.id = $1 
							AND p.deleted_at IS NULL 
							AND pp.deleted_at IS NULL)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, id).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r ProductRepoImpl) GetMasterProducts(c context.Context, queryParams queryparams.QueryParams) ([]entity.ProductDetail, error) {
	products := []entity.ProductDetail{}
	querIndex := 1

	query := `SELECT p.id, p.name
				FROM products p 
				WHERE p.is_active = True AND p.deleted_at IS NULL`

	var params []any

	query += queryparams.AddMasterConditionQuery(&params, queryParams, &querIndex)
	query += queryparams.AddMasterSortByQuery(queryParams)
	query += queryparams.AddMasterPaginationQuery(&params, queryParams, &querIndex)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.ProductDetail
		err := rows.Scan(
			&product.ID,
			&product.Name,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r ProductRepoImpl) GetTotalMasterProduct(c context.Context, queryParams queryparams.QueryParams) (int, error) {
	queryParams.Limit = 0
	queryParams.Page = 0
	queryParams.SortBy = ""
	queryParams.Order = ""
	querIndex := 1

	query := `SELECT COUNT(*)
				FROM products p 
				WHERE p.is_active = True AND p.deleted_at IS NULL`

	var params []any

	query += queryparams.AddMasterConditionQuery(&params, queryParams, &querIndex)
	query += queryparams.AddMasterPaginationQuery(&params, queryParams, &querIndex)

	var totalProduct int
	err := r.db.QueryRow(query, params...).Scan(&totalProduct)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return totalProduct, nil
}

func (r ProductRepoImpl) GetCategoryBoundary(c context.Context) (*entity.CategoryBoundary, error) {
	categoryBoundary := entity.CategoryBoundary{}

	query := `SELECT MAX(id), MIN(id)
				FROM product_categories pc 
				WHERE deleted_at IS NULL`

	err := r.db.QueryRow(query).Scan(&categoryBoundary.Maximum, &categoryBoundary.Minimum)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return &categoryBoundary, nil
}
