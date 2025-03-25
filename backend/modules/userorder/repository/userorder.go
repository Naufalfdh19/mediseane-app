package repository

import (
	"context"
	"database/sql"
	"fmt"
	"montelukast/modules/userorder/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"
)

type UserOrderRepo interface {
	IsOrderExistsByID(c context.Context, orderID int, userID int) (bool, error)
	UpdatePaymentStatus(c context.Context, orderDetailID int) error
	GetOrderDetailIDByID(c context.Context, orderID int) ([]int, error)
	GetStatusesByID(c context.Context, orderID int) ([]string, error)
	IsOrderDetailExists(c context.Context, order entity.Order) (bool, error)
	GetStatusByOrderDetailID(c context.Context, orderDetailID int) (string, error)
	UpdateDeliveryStatusByOrderDetailID(c context.Context, orderDetailID int) error
	GetDetailedOrdersByUserID(c context.Context, filter entity.OrderFilter) ([]entity.OrderProductDetail, error)
}

type userOrderRepoImpl struct {
	db *sql.DB
}

func NewUserOrderRepo(dbConn *sql.DB) userOrderRepoImpl {
	return userOrderRepoImpl{
		db: dbConn,
	}
}

func (r userOrderRepoImpl) IsOrderExistsByID(c context.Context, orderID int, userID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM orders WHERE id = $2 AND user_id = $1 AND  deleted_at IS NULL) FOR UPDATE `

	var exists bool
	err := r.db.QueryRow(query, userID, orderID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r userOrderRepoImpl) IsOrderDetailExists(c context.Context, order entity.Order) (bool, error) {
	query := `SELECT EXISTS (
				SELECT 1 
				FROM order_details od
				JOIN orders o 
				ON od.order_id = o.id 
				WHERE o.user_id = $1 AND od.id = $2 AND od.deleted_at IS NULL
			)`

	var exists bool
	err := r.db.QueryRowContext(c, query, order.UserID, order.OrderDetails[0].ID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r userOrderRepoImpl) GetOrderDetailIDByID(c context.Context, orderID int) ([]int, error) {
	tx := transaction.ExtractTx(c)

	orderDetailIDs := []int{}

	query := `SELECT od.id
				FROM order_details od 
				WHERE order_id = $1 `

	var err error
	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.QueryContext(c, query, orderID)
	} else {
		rows, err = r.db.QueryContext(c, query, orderID)
	}

	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var orderDetailID int
		err := rows.Scan(
			&orderDetailID,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		orderDetailIDs = append(orderDetailIDs, orderDetailID)
	}

	return orderDetailIDs, nil
}

func (r userOrderRepoImpl) GetStatusesByID(c context.Context, orderID int) ([]string, error) {
	tx := transaction.ExtractTx(c)

	allStatus := []string{}

	query := `SELECT od.status
				FROM order_details od
				WHERE order_id = $1 `

	var err error
	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.QueryContext(c, query, orderID)
	} else {
		rows, err = r.db.QueryContext(c, query, orderID)
	}

	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		err := rows.Scan(
			&status,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		allStatus = append(allStatus, status)
	}

	return allStatus, nil
}

func (r userOrderRepoImpl) UpdatePaymentStatus(c context.Context, orderDetailID int) error {
	tx := transaction.ExtractTx(c)

	query := `UPDATE order_details
				SET status = $2, updated_at = NOW() 
				WHERE id = $1 AND deleted_at IS NULL`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(c, query, orderDetailID, appconstant.StatusProcessing)
	} else {
		_, err = r.db.ExecContext(c, query, orderDetailID, appconstant.StatusProcessing)
	}
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r userOrderRepoImpl) UpdateDeliveryStatusByOrderDetailID(c context.Context, orderDetailID int) error {
	query := `UPDATE order_details 
				SET status = $2, updated_at = NOW() 
				WHERE id = $1 AND deleted_at IS NULL`

	var err error
	_, err = r.db.ExecContext(c, query, orderDetailID, appconstant.StatusDelivered)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return nil
}

func (r userOrderRepoImpl) GetStatusByOrderDetailID(c context.Context, orderDetailID int) (string, error) {
	query := `SELECT status FROM order_details WHERE id = $1 AND deleted_at IS NULL`

	var status string
	err := r.db.QueryRowContext(c, query, orderDetailID).Scan(&status)
	if err != nil {
		return "", apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return status, nil
}

func (r userOrderRepoImpl) GetDetailedOrdersByUserID(c context.Context, filter entity.OrderFilter) ([]entity.OrderProductDetail, error) {
	orderProductDetails := []entity.OrderProductDetail{}
	query := `SELECT 
				o.id order_id, o.total_price, o.created_at, 
				od.id order_detail_id, od.pharmacy_id, od.status, od.logistic_price,
				opd.id order_detail_product_id, opd.pharmacy_product_id, opd.quantity, opd.price subtotal,
				p.name, p.manufacture, p.image[1],
				pmc.name pharmacy_name
			FROM
				orders o
			JOIN 
				order_details od ON o.id = od.order_id
			JOIN 
				order_product_details opd ON od.id = opd.order_detail_id
			JOIN
				pharmacy_products pp ON pp.id = opd.pharmacy_product_id
			JOIN
				products p ON p.id = pp.product_id 
			JOIN 
				pharmacies pmc ON pmc.id = od.pharmacy_id
			WHERE 
				`

	paramQuery, sortPage, args := OrderParam(filter)
	userFilter := fmt.Sprintf("o.user_id = %d", filter.UserID)
	finalQuery := query + userFilter + paramQuery + sortPage
	rows, err := r.db.QueryContext(c, finalQuery, args...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var orderProductDetail entity.OrderProductDetail
		err := rows.Scan(
			&orderProductDetail.OrderDetail.Order.ID,
			&orderProductDetail.OrderDetail.Order.TotalPrice,
			&orderProductDetail.OrderDetail.Order.CreatedAt,
			&orderProductDetail.OrderDetail.ID,
			&orderProductDetail.OrderDetail.PharmacyID,
			&orderProductDetail.OrderDetail.Status,
			&orderProductDetail.OrderDetail.LogisticPrice,
			&orderProductDetail.ID,
			&orderProductDetail.PharmacyProductID,
			&orderProductDetail.Quantity,
			&orderProductDetail.Subtotal,
			&orderProductDetail.Name,
			&orderProductDetail.Manufacturer,
			&orderProductDetail.Image,
			&orderProductDetail.PharmacyName,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		orderProductDetails = append(orderProductDetails, orderProductDetail)
	}

	err = rows.Err()
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return orderProductDetails, nil
}
