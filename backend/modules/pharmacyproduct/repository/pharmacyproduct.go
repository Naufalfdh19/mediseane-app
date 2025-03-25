package repository

import (
	"context"
	"database/sql"
	"montelukast/modules/pharmacyproduct/entity"
	"montelukast/modules/pharmacyproduct/queryparams"
	productEntity "montelukast/modules/product/entity"
	appconstant "montelukast/pkg/constant"
	"montelukast/pkg/dateconverter"
	apperror "montelukast/pkg/error"
	"time"

	"github.com/go-redis/redis/v8"
)

type PharmacyProductRepo interface {
	IsPharmacyProductExistsByID(c context.Context, id int) (bool, error)
	IsPharmacyProductExistsByIDAndPharmacy(c context.Context, pharmacyProductID int, pharmacyID int) (bool, error)
	IsPharmacyProductExists(c context.Context, pharmacyID, productID int) (bool, error)
	IsProductExistsByID(c context.Context, id int) (bool, error)
	GetStockByID(c context.Context, id int) (int, error)
	AddPharmacyProduct(c context.Context, pharmacyProduct entity.PharmacyProduct) error
	GetPharmacyIDbyPharmacistID(c context.Context, pharmacistID int) (int, error)
	UpdatePharmacyProduct(c context.Context, pharmacistProduct entity.PharmacyProduct) error
	GetPharmacyIDbyPharmacyProductID(c context.Context, pharmacyProductID int) (int, error)
	GetStockUpdatedDate(c context.Context) (string, error)
	SetStockUpdatedDate(c context.Context) error
	DeletePharmacyProduct(c context.Context, pharmacyProductID int) error
	GetPharmacyProduct(c context.Context, productPharmacyID int, pharmacyID int) (*productEntity.ProductDetail, error)
	GetPharmacyProducts(c context.Context, queryParams queryparams.QueryParams, pharmacyID int) ([]productEntity.ProductDetail, error)
	GetTotalProduct(c context.Context, queryParams queryparams.QueryParams, pharmacyID int) (int, error)
}

type pharmacyProductRepoImpl struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewPharmacyProductRepo(dbConn *sql.DB, rdbConn *redis.Client) pharmacyProductRepoImpl {
	return pharmacyProductRepoImpl{
		db:  dbConn,
		rdb: rdbConn,
	}
}

func (r pharmacyProductRepoImpl) IsPharmacyProductExistsByID(c context.Context, id int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacy_products WHERE id = $1 AND deleted_at IS NULL)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, id).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r pharmacyProductRepoImpl) GetStockByID(c context.Context, id int) (int, error) {
	query := `SELECT stock FROM pharmacy_products WHERE id = $1 AND deleted_at IS NULL`

	var stock int
	err := r.db.QueryRowContext(c, query, id).Scan(&stock)
	if err != nil {
		return -1, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return stock, nil
}

func (r pharmacyProductRepoImpl) IsPharmacyProductExistsByIDAndPharmacy(c context.Context, pharmacyProductID int, pharmacyID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacy_products WHERE id = $1 AND pharmacy_id = $2 AND deleted_at IS NULL)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, pharmacyProductID, pharmacyID).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r pharmacyProductRepoImpl) IsPharmacyProductExists(c context.Context, pharmacyID, productID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacy_products WHERE pharmacy_id = $1 AND product_id = $2 AND deleted_at IS NULL)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, pharmacyID, productID).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r pharmacyProductRepoImpl) IsProductExistsByID(c context.Context, id int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM products WHERE id = $1 AND deleted_at IS NULL)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, id).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r pharmacyProductRepoImpl) AddPharmacyProduct(c context.Context, pharmacyProduct entity.PharmacyProduct) error {

	query := `INSERT INTO pharmacy_products (pharmacy_id, product_id, stock, price, is_active)
			  VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query,
		pharmacyProduct.PharmacyID,
		pharmacyProduct.ProductID,
		pharmacyProduct.Stock,
		pharmacyProduct.Price,
		pharmacyProduct.IsActive,
	)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r pharmacyProductRepoImpl) GetPharmacyIDbyPharmacistID(c context.Context, pharmacistID int) (int, error) {
	query := `SELECT pharmacy_id
				FROM pharmacist_details
				WHERE pharmacist_id = $1 AND deleted_at IS NULL`

	var pharmacyID int
	err := r.db.QueryRow(query, pharmacistID).Scan(&pharmacyID)
	if err != nil && err != sql.ErrNoRows {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return pharmacyID, nil
}

func (r pharmacyProductRepoImpl) GetPharmacyIDbyPharmacyProductID(c context.Context, pharmacyProductID int) (int, error) {
	query := `SELECT pharmacy_id
				FROM pharmacy_products
				WHERE id = $1 AND deleted_at IS NULL`

	var pharmacyID int
	err := r.db.QueryRow(query, pharmacyProductID).Scan(&pharmacyID)
	if err != nil && err != sql.ErrNoRows {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return pharmacyID, nil
}

func (r pharmacyProductRepoImpl) UpdatePharmacyProduct(c context.Context, pharmacistProduct entity.PharmacyProduct) error {
	query := `UPDATE pharmacy_products
				SET stock = $2, is_active = $3, updated_at = NOW()
				WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(query, pharmacistProduct.ID, pharmacistProduct.Stock, pharmacistProduct.IsActive)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r pharmacyProductRepoImpl) GetStockUpdatedDate(c context.Context) (string, error) {
	res := r.rdb.Get(c, "stock-update-time")
	err := res.Err()
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	stockUpdatedTime, err := res.Result()
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return stockUpdatedTime, nil
}

func (r pharmacyProductRepoImpl) SetStockUpdatedDate(c context.Context) error {
	currentDate := dateconverter.GetCurrentDate()
	ttl := time.Duration(10) * time.Second

	opt := r.rdb.Set(c, "stock-update-time", currentDate, ttl)
	err := opt.Err()
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r pharmacyProductRepoImpl) DeletePharmacyProduct(c context.Context, pharmacyProductID int) error {
	query := `UPDATE pharmacy_products
				SET deleted_at = NOW()
				WHERE id = $1`

	_, err := r.db.Exec(query, pharmacyProductID)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r pharmacyProductRepoImpl) GetPharmacyProduct(c context.Context, productPharmacyID int, pharmacyID int) (*productEntity.ProductDetail, error) {
	query := `SELECT p.id, pp.id, p.name, p.generic_name, p.manufacture, pc.name, pf.name, pp.stock, pp.is_active, pp.created_at, pp.price
				FROM products p 
				JOIN pharmacy_products pp ON pp.id = $1 AND pp.pharmacy_id = $2 AND pp.product_id = p.id AND pp.deleted_at IS NULL
				JOIN product_classifications pc ON pc.id = p.product_classification_id AND pc.deleted_at IS NULL
				JOIN product_forms pf ON pf.id = p.product_form_id AND pf.deleted_at IS NULL
				WHERE p.deleted_at IS NULL;`

	var product productEntity.ProductDetail
	err := r.db.QueryRowContext(c, query, productPharmacyID, pharmacyID).Scan(
		&product.ID,
		&product.PharmacyProductID,
		&product.Name,
		&product.GenericName,
		&product.Manufacture,
		&product.ProductClassification,
		&product.ProductForm,
		&product.Stock,
		&product.IsActive,
		&product.CreatedAt,
		&product.Price,
	)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return &product, nil
}

func (r pharmacyProductRepoImpl) GetPharmacyProducts(c context.Context, queryParams queryparams.QueryParams, pharmacyID int) ([]productEntity.ProductDetail, error) {
	products := []productEntity.ProductDetail{}
	querIndex := 1

	query := `SELECT p.id, pp.id, p.name, p.generic_name, p.manufacture, pc.name, pf.name, pp.stock, pp.is_active, pp.created_at, pp.price
				FROM products p 
				JOIN pharmacy_products pp ON pp.product_id = p.id AND pp.deleted_at IS NULL
				JOIN product_classifications pc ON pc.id = p.product_classification_id AND pc.deleted_at IS NULL
				JOIN product_forms pf ON pf.id = p.product_form_id AND pf.deleted_at IS NULL
				WHERE pp.pharmacy_id = $1 AND p.deleted_at IS NULL`

	var params []any
	params = append(params, pharmacyID)
	querIndex++

	query += queryparams.AddQueryParams(&params, queryParams, &querIndex)
	query += queryparams.AddSortByQuery(queryParams)
	query += queryparams.AddPaginationQuery(&params, queryParams, &querIndex)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var product productEntity.ProductDetail
		err := rows.Scan(
			&product.ID,
			&product.PharmacyProductID,
			&product.Name,
			&product.GenericName,
			&product.Manufacture,
			&product.ProductClassification,
			&product.ProductForm,
			&product.Stock,
			&product.IsActive,
			&product.CreatedAt,
			&product.Price,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r pharmacyProductRepoImpl) GetTotalProduct(c context.Context, queryParams queryparams.QueryParams, pharmacyID int) (int, error) {
	queryParams.Limit = 0
	queryParams.Page = 0
	queryParams.SortBy = ""
	queryParams.Order = ""
	querIndex := 1

	query := `SELECT COUNT(*)
				FROM products p 
				JOIN pharmacy_products pp ON pp.product_id = p.id AND pp.deleted_at IS NULL
				JOIN product_classifications pc ON pc.id = p.product_classification_id AND pc.deleted_at IS NULL
				JOIN product_forms pf ON pf.id = p.product_form_id AND pf.deleted_at IS NULL
				WHERE pp.pharmacy_id = $1 AND p.deleted_at IS NULL`

	var params []any
	params = append(params, pharmacyID)
	querIndex++

	query += queryparams.AddQueryParams(&params, queryParams, &querIndex)
	query += queryparams.AddPaginationQuery(&params, queryParams, &querIndex)

	var totalProduct int
	err := r.db.QueryRow(query, params...).Scan(&totalProduct)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return totalProduct, nil
}
