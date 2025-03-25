package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"montelukast/modules/checkout/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
)

type CheckoutRepo interface {
	AddOrder(c context.Context, totalPrice decimal.Decimal, userID int) (int, error)
	GetCheckoutCartRedis(c context.Context, cartID string, userID int) (result *entity.ListGroupedCartItem, err error)
	DeleteCartRedis(c context.Context, cartID string, userID int) error
	AddCheckoutOrderDetail(c context.Context, listData []entity.DeliveryPriceData, orderID int) ([]int, error)
	AddOrderProductDetails(c context.Context, orderDetailID int, products entity.GroupedCartItem) error
	UpdateStock(c context.Context, product entity.CartItem) error
	GetQuantity(c context.Context, pharmacyProductIDint int) (int, error)
	IsActive(c context.Context, pharmacyProductIDint int) (bool, error)
	IsCartItemExistsByIDAndUserID(c context.Context, cartItem entity.CartItem) (bool, error)
	DeleteCartItemByID(c context.Context, id int) error
	IsOrderFromUser(c context.Context, order_id int, user_id int) (bool, error)
	IsOrderExistByID(c context.Context, orderID int, userID int) (bool, error)
	CancelOrder(c context.Context, orderID int) error
}

type checkOutRepoImpl struct {
	db      *sql.DB
	redisDB *redis.Client
}

func NewCheckoutRepo(db *sql.DB, redisDB *redis.Client) checkOutRepoImpl {
	return checkOutRepoImpl{
		db:      db,
		redisDB: redisDB,
	}
}

func (r *checkOutRepoImpl) IsCartItemExistsByIDAndUserID(c context.Context, cartItem entity.CartItem) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM carts WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL)`
	var isExists bool
	err := r.db.QueryRowContext(c, query, cartItem.ID, cartItem.UserID).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r *checkOutRepoImpl) DeleteCartItemByID(c context.Context, id int) error {
	tx := transaction.ExtractTx(c)
	query := `UPDATE 
				carts
			SET 
				deleted_at = NOW() 
			WHERE 
				id = $1`

	_, err := tx.ExecContext(c, query, id)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r *checkOutRepoImpl) IsActive(c context.Context, pharmacyProductIDint int) (bool, error) {
	var is_active bool
	var err error
	tx := transaction.ExtractTx(c)
	query := `SELECT is_active
				FROM pharmacy_products 
				WHERE id = $1`
	if tx != nil {
		err = tx.QueryRowContext(c, query, pharmacyProductIDint).Scan(&is_active)
	} else {
		err = r.db.QueryRowContext(c, query, pharmacyProductIDint).Scan(&is_active)
	}
	if err != nil {
		return is_active, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return is_active, nil
}

func (r *checkOutRepoImpl) UpdateStock(c context.Context, product entity.CartItem) error {
	tx := transaction.ExtractTx(c)
	query := `UPDATE pharmacy_products 
				SET stock = stock - $2
				WHERE id = $1 AND is_active IS TRUE`
	_, err := tx.ExecContext(c, query, product.PharmacyProductID, product.Quantity)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r *checkOutRepoImpl) GetQuantity(c context.Context, pharmacyProductIDint int) (int, error) {
	query := `SELECT stock FROM pharmacy_products
	WHERE id = $1`
	var quantity int
	err := r.db.QueryRowContext(c, query, pharmacyProductIDint).Scan(&quantity)
	if err != nil {
		return -1, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return quantity, nil
}

func (r *checkOutRepoImpl) AddOrderProductDetails(c context.Context, orderDetailID int, products entity.GroupedCartItem) error {
	tx := transaction.ExtractTx(c)
	valueStrings := make([]string, 0, len(products.Items))
	valueArgs := make([]interface{}, 0, len(products.Items))
	for i, product := range products.Items {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, orderDetailID, product.PharmacyProductID, product.Quantity, product.Subtotal)
	}
	query := fmt.Sprintf(`INSERT INTO order_product_details(
	order_detail_id, pharmacy_product_id, quantity, price)
	VALUES %s`, strings.Join(valueStrings, ","))
	_, err := tx.ExecContext(c, query, valueArgs...)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r *checkOutRepoImpl) AddCheckoutOrderDetail(c context.Context, listData []entity.DeliveryPriceData, orderID int) ([]int, error) {
	tx := transaction.ExtractTx(c)
	valueStrings := make([]string, 0, len(listData))
	valueArgs := make([]interface{}, 0, len(listData))
	for i, data := range listData {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, orderID, data.PharmacyID, data.LogisticPrice, appconstant.StatusPending)
	}
	query := fmt.Sprintf(`INSERT INTO order_details(
	order_id,pharmacy_id, logistic_price, status)
	VALUES %s RETURNING id`, strings.Join(valueStrings, ","))
	rows, err := tx.QueryContext(c, query, valueArgs...)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *checkOutRepoImpl) AddOrder(c context.Context, totalPrice decimal.Decimal, userID int) (int, error) {
	tx := transaction.ExtractTx(c)
	var id int
	query := `INSERT INTO orders 
	(user_id, total_price) 
	VALUES ($1,$2) RETURNING id`
	err := tx.QueryRowContext(c, query, userID, totalPrice).Scan(&id)
	if err != nil {
		return -1, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return id, nil
}

func (r *checkOutRepoImpl) GetCheckoutCartRedis(c context.Context, cartID string, userID int) (result *entity.ListGroupedCartItem, err error) {
	key := fmt.Sprintf("checkout:%d:cartIds:%s", userID, cartID)
	serializedCarts, err := r.redisDB.Get(c, key).Result()
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	var carts entity.ListGroupedCartItem
	err = json.Unmarshal([]byte(serializedCarts), &carts)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return &carts, nil
}

func (r *checkOutRepoImpl) DeleteCartRedis(c context.Context, cartID string, userID int) error {
	key := fmt.Sprintf("checkout:%d:cartIds:%s", userID, cartID)
	_, err := r.redisDB.Del(c, key).Result()
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r *checkOutRepoImpl) IsOrderFromUser(c context.Context, order_id int, user_id int) (bool, error) {
	var orderFromUser bool
	query := `SELECT EXISTS (SELECT 1 FROM orders 
				WHERE user_id = $1 AND id = $2)`
	err := r.db.QueryRowContext(c, query, user_id, order_id).Scan(&orderFromUser)
	if err != nil {
		return orderFromUser, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return orderFromUser, nil
}

func (r *checkOutRepoImpl) IsOrderExistByID(c context.Context, orderID int, userID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM orders WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL)`
	var exists bool
	err := r.db.QueryRowContext(c, query, orderID, userID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return exists, nil
}

func (r *checkOutRepoImpl) CancelOrder(c context.Context, orderID int) error {
	query := `UPDATE order_details
				SET 
					status = $2,
					updated_at = NOW()
				WHERE order_id = $1`
	_, err := r.db.ExecContext(c, query, orderID, appconstant.StatusCancelled)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}
