package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"montelukast/modules/cart/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
)

type CartRepo interface {
	IsCartItemExistsByProductIDAndUserID(c context.Context, cartItem entity.CartItem) (bool, error)
	IsCartItemExistsByIDAndUserID(c context.Context, cartItem entity.CartItem) (bool, error)
	AddCartItem(c context.Context, cartItem entity.CartItem) error
	UpdateCartItemQuantity(c context.Context, cartItem entity.CartItem) error
	DeleteCartItemByID(c context.Context, id int) error
	GetDetailedCartItemsByUserID(c context.Context, userID int) ([]entity.CartItem, error)
	GetCartItemsByUserID(c context.Context, userID int) ([]entity.CartItem, error)
	GetSelectedCartItems(c context.Context, userID int, ids []int) ([]entity.CartItem, error)
	GetQuantityByProductIDAndUserID(c context.Context, cartItem entity.CartItem) (int, error)
	AddCheckoutItem(c context.Context, carts entity.ListGroupedCartItem, userID int) (err error)
	GetCheckoutCartRedis(c context.Context, cartID string, userID int) (result *entity.ListGroupedCartItem, err error)
	IsUserVerified(c context.Context, userID int) (bool, error)
}

type cartRepoImpl struct {
	db      *sql.DB
	redisDB *redis.Client
}

func NewCartRepo(dbConn *sql.DB, redisDB *redis.Client) cartRepoImpl {
	return cartRepoImpl{
		db:      dbConn,
		redisDB: redisDB,
	}
}

func (r cartRepoImpl) GetCheckoutCartRedis(c context.Context, cartID string, userID int) (result *entity.ListGroupedCartItem, err error) {
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

func (r cartRepoImpl) AddCheckoutItem(c context.Context, carts entity.ListGroupedCartItem, userID int) (err error) {
	data, err := json.Marshal(carts)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	key := fmt.Sprintf("checkout:%d:cartIds:%s", userID, carts.ID)
	_, err = r.redisDB.Set(c, key, data, appconstant.CartRedisExpiration).Result()
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r cartRepoImpl) IsCartItemExistsByProductIDAndUserID(c context.Context, cartItem entity.CartItem) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM carts WHERE user_id = $1 AND pharmacy_product_id = $2 AND deleted_at IS NULL)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, cartItem.UserID, cartItem.PharmacyProductID).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r cartRepoImpl) IsCartItemExistsByIDAndUserID(c context.Context, cartItem entity.CartItem) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM carts WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, cartItem.ID, cartItem.UserID).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}

func (r cartRepoImpl) GetQuantityByProductIDAndUserID(c context.Context, cartItem entity.CartItem) (int, error) {
	query := `SELECT 
				c.quantity 
			FROM 
				carts c 
			WHERE
				c.user_id = $1 AND c.pharmacy_product_id = $2 AND c.deleted_at IS NULL;`

	var quantity int
	err := r.db.QueryRowContext(c, query, cartItem.UserID, cartItem.PharmacyProductID).Scan(&quantity)
	if err != nil && err != sql.ErrNoRows {
		return -1, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return quantity, nil
}

func (r cartRepoImpl) AddCartItem(c context.Context, cartItem entity.CartItem) error {
	query := `INSERT INTO carts (user_id, pharmacy_product_id, quantity)
			  VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(c, query, cartItem.UserID, cartItem.PharmacyProductID, cartItem.Quantity)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r cartRepoImpl) UpdateCartItemQuantity(c context.Context, cartItem entity.CartItem) error {
	query := `UPDATE 
				carts
			SET 
				quantity = quantity + $3,
				updated_at = NOW() 
			WHERE 
				user_id = $1 AND pharmacy_product_id = $2`

	_, err := r.db.ExecContext(c, query, cartItem.UserID, cartItem.PharmacyProductID, cartItem.Quantity)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r cartRepoImpl) DeleteCartItemByID(c context.Context, id int) error {
	query := `UPDATE 
				carts
			SET 
				deleted_at = NOW() 
			WHERE 
				id = $1`

	_, err := r.db.ExecContext(c, query, id)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return nil
}

func (r cartRepoImpl) GetDetailedCartItemsByUserID(c context.Context, userID int) ([]entity.CartItem, error) {
	cartItems := []entity.CartItem{}
	query := `SELECT 
				c.id, c.user_id, c.pharmacy_product_id, pp.pharmacy_id, pmc.name, c.quantity, pp.price, p.image, p.name, p.manufacture
			FROM 
				carts c 
			JOIN 
				pharmacy_products pp ON c.user_id = $1 AND c.pharmacy_product_id = pp.id AND c.deleted_at IS NULL
			JOIN 
				products p ON pp.product_id = p.id
			JOIN
				pharmacies pmc ON pp.pharmacy_id = pmc.id AND pmc.deleted_at IS NULL
			ORDER BY 
				pp.pharmacy_id, c.created_at;`

	rows, err := r.db.QueryContext(c, query, userID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cartItem entity.CartItem
		var price string
		err := rows.Scan(
			&cartItem.ID,
			&cartItem.UserID,
			&cartItem.PharmacyProductID,
			&cartItem.PharmacyID,
			&cartItem.PharmacyName,
			&cartItem.Quantity,
			&price,
			&cartItem.Image,
			&cartItem.Name,
			&cartItem.Manufacturer,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		priceNum, err := decimal.NewFromString(price)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrGetCart, apperror.ErrConvertVariableType, err)
		}
		cartItem.Subtotal = priceNum.Mul(decimal.NewFromInt32(int32(cartItem.Quantity)))
		cartItems = append(cartItems, cartItem)
	}

	err = rows.Err()
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return cartItems, nil
}

func (r cartRepoImpl) GetCartItemsByUserID(c context.Context, userID int) ([]entity.CartItem, error) {
	cartItems := []entity.CartItem{}
	query := `SELECT 
				c.id, c.pharmacy_product_id, c.quantity 
			FROM 
				carts c 
			JOIN 
				pharmacy_products pp ON c.user_id = $1 AND c.pharmacy_product_id = pp.id AND c.deleted_at IS NULL
			JOIN 
				products p ON pp.product_id = p.id
			JOIN
				pharmacies pmc ON pp.pharmacy_id = pmc.id AND pmc.deleted_at IS NULL
			WHERE
				c.user_id = $1 AND c.deleted_at IS NULL;`

	rows, err := r.db.QueryContext(c, query, userID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()
	for rows.Next() {
		var cartItem entity.CartItem
		err := rows.Scan(
			&cartItem.ID,
			&cartItem.PharmacyProductID,
			&cartItem.Quantity,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		cartItems = append(cartItems, cartItem)
	}

	err = rows.Err()
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return cartItems, nil
}

func (r cartRepoImpl) GetSelectedCartItems(c context.Context, userID int, ids []int) ([]entity.CartItem, error) {
	cartItems := []entity.CartItem{}

	buf := bytes.NewBufferString(`SELECT 
									c.id, c.user_id, c.pharmacy_product_id, pp.pharmacy_id, pmc.name, c.quantity, pp.price, p.image, p.name, p.manufacture
								FROM 
									carts c 
								JOIN 
									pharmacy_products pp ON c.user_id = $1 AND c.id IN(`)

	for i, id := range ids {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(id))
	}

	buf.WriteString(`) AND c.pharmacy_product_id = pp.id AND c.deleted_at IS NULL
					JOIN 
						products p ON pp.product_id = p.id
					JOIN
						pharmacies pmc ON pp.pharmacy_id = pmc.id;`)

	rows, err := r.db.QueryContext(c, buf.String(), userID)
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cartItem entity.CartItem
		var price string
		err := rows.Scan(
			&cartItem.ID,
			&cartItem.UserID,
			&cartItem.PharmacyProductID,
			&cartItem.PharmacyID,
			&cartItem.PharmacyName,
			&cartItem.Quantity,
			&price,
			&cartItem.Image,
			&cartItem.Name,
			&cartItem.Manufacturer,
		)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
		}
		priceNum, err := decimal.NewFromString(price)
		if err != nil {
			return nil, apperror.NewErrInternalServerError(appconstant.FieldErrGetCart, apperror.ErrConvertVariableType, err)
		}
		cartItem.Subtotal = priceNum.Mul(decimal.NewFromInt32(int32(cartItem.Quantity)))
		cartItems = append(cartItems, cartItem)
	}

	err = rows.Err()
	if err != nil {
		return nil, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}

	return cartItems, nil
}

func (r cartRepoImpl) IsUserVerified(c context.Context, userID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1 AND is_verified = True AND deleted_at IS NULL)`

	var isExists bool
	err := r.db.QueryRowContext(c, query, userID).Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return isExists, apperror.NewErrInternalServerError(appconstant.FieldErrServer, apperror.ErrInternalServer, err)
	}
	return isExists, nil
}