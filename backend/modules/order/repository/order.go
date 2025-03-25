package repository

import (
	"context"
	"database/sql"
	"montelukast/modules/order/entity"
	queryparams "montelukast/modules/order/query_params"
	productEntity "montelukast/modules/product/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"
)

type OrderRepo interface {
	GetTotalOrder(c context.Context, queryParams queryparams.QueryParams, userID int) (int, error)
	GetOrders(c context.Context, queryParams queryparams.QueryParams, addressID int) ([]entity.OrderDetail, error)
	IsPharmacistExistsByID(c context.Context, pharmacistID int) (bool, error)
	GetPharmacyIDByPharmacistID(c context.Context, pharmacistID int) (int, error)
	IsOrderDetailExistsByID(c context.Context, orderDetailID int) (bool, error)
	IsPharmacyExistsByPharmacistID(c context.Context, pharmacistID int) (bool, error)
	GetOrderDetailByID(c context.Context, orderDetailID int) (*entity.OrderDetail, error)
	GetOrderedProduct(c context.Context, orderDetail int, pharmacyID int) ([]productEntity.ProductDetail, error)
	GetPharmacyIDByOrderID(c context.Context, orderDetailID int) (int, error)
	DeleteOrderDetails(c context.Context, orderDetailID int) error
	DeleteOrderProductDetails(c context.Context, orderDetailID int) error
	GetOrderStatusByID(c context.Context, orderDetailID int) (string, error)
	UpdateOrderStatus(c context.Context, orderDetailID int) error
	GetOrderedProductsQuantity(c context.Context, orderDetailID int) (map[int]int, error)
	UpdateProductQuantity(c context.Context, pharmacyProductID, quantity int) error
	UpdateOrderStatusDelivered(c context.Context, orderDetailID int) error
}

type orderRepoImpl struct {
	db *sql.DB
}

func NewOrderRepo(dbConn *sql.DB) orderRepoImpl {
	return orderRepoImpl{
		db: dbConn,
	}
}

func (r orderRepoImpl) GetOrders(c context.Context, queryParams queryparams.QueryParams, pharmacyID int) ([]entity.OrderDetail, error) {
	orders := []entity.OrderDetail{}

	query := `SELECT od.id, od.status, o.created_at 
				FROM orders o
				JOIN order_details od ON od.order_id = o.id
				WHERE pharmacy_id = $1 AND o.deleted_at IS NULL`

	var params []any
	params = append(params, pharmacyID)
	query += queryparams.AddQueryParams(&params, queryParams)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.OrderDetail
		err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (r orderRepoImpl) GetTotalOrder(c context.Context, queryParams queryparams.QueryParams, pharmacyID int) (int, error) {
	queryParams.Limit = 0
	queryParams.Page = 0
	queryParams.SortBy = ""
	queryParams.Order = ""

	query := `SELECT COUNT(*)
				FROM orders o
				JOIN order_details od ON od.order_id = o.id
				WHERE pharmacy_id = $1 AND o.deleted_at IS NULL`

	var params []any
	params = append(params, pharmacyID)
	query += queryparams.AddQueryParams(&params, queryParams)

	var totalOrder int
	err := r.db.QueryRow(query, params...).Scan(&totalOrder)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return totalOrder, nil
}

func (r orderRepoImpl) IsPharmacistExistsByID(c context.Context, pharmacistID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacist_details WHERE pharmacist_id = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, pharmacistID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r orderRepoImpl) GetPharmacyIDByPharmacistID(c context.Context, pharmacistID int) (int, error) {
	query := `SELECT pharmacy_id 
				FROM pharmacist_details pd 
				WHERE pharmacist_id = $1 AND pharmacy_id IS NOT NULL AND deleted_at IS NULL`

	var pharmacyID int
	err := r.db.QueryRow(query, pharmacistID).Scan(&pharmacyID)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return pharmacyID, nil
}

func (r orderRepoImpl) IsOrderDetailExistsByID(c context.Context, orderDetailID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM order_details WHERE id = $1 AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, orderDetailID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r orderRepoImpl) GetOrderDetailByID(c context.Context, orderDetailID int) (*entity.OrderDetail, error) {
	query := `select id, status, created_at
				from order_details od 
				where id = $1 AND deleted_at IS NULL`

	var orderDetail entity.OrderDetail
	err := r.db.QueryRow(query, orderDetailID).Scan(
		&orderDetail.ID,
		&orderDetail.Status,
		&orderDetail.CreatedAt,
	)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return &orderDetail, nil
}

func (r orderRepoImpl) IsPharmacyExistsByPharmacistID(c context.Context, pharmacistID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM pharmacist_details WHERE pharmacist_id = $1 AND pharmacy_id IS NOT NULL AND deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, pharmacistID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r orderRepoImpl) GetPharmacyIDByOrderID(c context.Context, orderDetailID int) (int, error) {
	query := `select pharmacy_id 
				from order_details pd 
				where id = $1 AND deleted_at IS NULL`

	var pharmacyID int
	err := r.db.QueryRow(query, orderDetailID).Scan(&pharmacyID)
	if err != nil {
		return 0, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return pharmacyID, nil
}

func (r orderRepoImpl) GetOrderedProduct(c context.Context, orderDetailID int, pharmacyID int) ([]productEntity.ProductDetail, error) {
	productOrders := []productEntity.ProductDetail{}

	query := `SELECT p.id, p.name, opd.quantity, p.image[1]
				FROM order_details od 
				JOIN order_product_details opd ON opd.order_detail_id = od.id
				JOIN pharmacy_products pp ON pp.id = opd.pharmacy_product_id 
				JOIN products p ON p.id = pp.product_id 
				WHERE od.id = $1 AND od.pharmacy_id = $2 AND od.deleted_at IS NULL AND opd.deleted_at IS NULL AND pp.deleted_at IS NULL AND p.deleted_at IS NULL`

	rows, err := r.db.Query(query, orderDetailID, pharmacyID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var productOrder productEntity.ProductDetail
		err := rows.Scan(
			&productOrder.ID,
			&productOrder.Name,
			&productOrder.Stock,
			&productOrder.Image,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		productOrders = append(productOrders, productOrder)
	}
	return productOrders, nil
}

func (r orderRepoImpl) GetOrderStatusByID(c context.Context, orderDetailID int) (string, error) {
	query := `select status
				from order_details pd 
				where id = $1 AND deleted_at IS NULL`

	var orderStatus string
	err := r.db.QueryRow(query, orderDetailID).Scan(&orderStatus)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return orderStatus, nil
}

func (r orderRepoImpl) DeleteOrderDetails(c context.Context, orderDetailID int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE order_details
				SET status = $2, deleted_at = NOW()
				WHERE id = $1 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, orderDetailID, appconstant.StatusCancelled)
	} else {
		_, err = r.db.ExecContext(c, query, orderDetailID, appconstant.StatusCancelled)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r orderRepoImpl) DeleteOrderProductDetails(c context.Context, orderDetailID int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE order_product_details
				SET deleted_at = NOW()
				WHERE order_detail_id = $1 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, orderDetailID)
	} else {
		_, err = r.db.ExecContext(c, query, orderDetailID)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r orderRepoImpl) UpdateOrderStatus(c context.Context, orderDetailID int) error {
	query := `UPDATE order_details
				SET status = $2, updated_at = NOW()
				WHERE id = $1`

	_, err := r.db.Exec(query, orderDetailID, appconstant.StatusShipped)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r orderRepoImpl) UpdateOrderStatusDelivered(c context.Context, orderDetailID int) error {
	query := `UPDATE order_details
				SET status = $2, updated_at = NOW()
				WHERE id = $1`

	_, err := r.db.Exec(query, orderDetailID, appconstant.StatusDelivered)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r orderRepoImpl) GetOrderedProductsQuantity(c context.Context, orderDetailID int) (map[int]int, error) {
	orderedProducts := make(map[int]int)

	query := `SELECT opd.pharmacy_product_id, opd.quantity 
				FROM order_details od 
				JOIN order_product_details opd ON opd.order_detail_id = od.id AND opd.deleted_at IS NULL
				WHERE od.id = $1 AND od.deleted_at IS NULL`

	rows, err := r.db.Query(query, orderDetailID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var pharmacyProductID int
		var quantity int
		err := rows.Scan(
			&pharmacyProductID,
			&quantity,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		orderedProducts[pharmacyProductID] = quantity
	}
	return orderedProducts, nil
}

func (r orderRepoImpl) UpdateProductQuantity(c context.Context, pharmacyProductID, quantity int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE pharmacy_products
				SET stock = stock + $2, updated_at = NOW()
				WHERE id = $1 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, pharmacyProductID, quantity)
	} else {
		_, err = r.db.ExecContext(c, query, pharmacyProductID, quantity)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}
