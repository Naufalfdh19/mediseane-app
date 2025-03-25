package usecase

import (
	"context"
	"montelukast/modules/cart/entity"
	"montelukast/modules/cart/repository"
	pharmacyproduct "montelukast/modules/pharmacyproduct/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"

	"github.com/google/uuid"
)

type CartUsecase interface {
	AddToCart(c context.Context, cartItem entity.CartItem) error
	DeleteFromCart(c context.Context, cartItem entity.CartItem) error
	GetGroupedCartItems(c context.Context, userID int) ([]entity.GroupedCartItem, error)
	GetCartItems(c context.Context, userID int) ([]entity.CartItem, error)
	GetSelectedCartItems(c context.Context, userID int, ids []int) (*entity.ListGroupedCartItem, error)
}

type cartUsecaseImpl struct {
	r  repository.CartRepo
	pp pharmacyproduct.PharmacyProductRepo
}

func NewCartUsecase(r repository.CartRepo, pp pharmacyproduct.PharmacyProductRepo) cartUsecaseImpl {
	return cartUsecaseImpl{
		r:  r,
		pp: pp,
	}
}

func (u cartUsecaseImpl) isStockSufficient(c context.Context, cartItem entity.CartItem, stock int) (bool, error) {
	quantity, err := u.r.GetQuantityByProductIDAndUserID(c, cartItem)
	if err != nil {
		return false, err
	}
	return (quantity + cartItem.Quantity) <= stock, nil
}

func (u cartUsecaseImpl) AddToCart(c context.Context, cartItem entity.CartItem) error {
	isVerified, err := u.r.IsUserVerified(c, cartItem.UserID)
	if err != nil {
		return err
	}
	if !isVerified {
		return apperror.NewErrStatusUnauthorized(appconstant.FieldErrAddToCart, apperror.ErrUserNotVerified, apperror.ErrUserNotVerified)
	}

	isProductExists, err := u.pp.IsPharmacyProductExistsByID(c, cartItem.PharmacyProductID)
	if err != nil {
		return err
	}
	if !isProductExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrAddToCart, apperror.ErrPharmacyProductNotExists, apperror.ErrPharmacyProductNotExists)
	}

	stock, err := u.pp.GetStockByID(c, cartItem.PharmacyProductID)
	if err != nil {
		return err
	}
	isAvailable, err := u.isStockSufficient(c, cartItem, stock)
	if err != nil {
		return err
	}
	if !isAvailable {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrAddToCart, apperror.ErrStockUnavailable, apperror.ErrStockUnavailable)
	}

	isAlreadyExists, err := u.r.IsCartItemExistsByProductIDAndUserID(c, cartItem)
	if err != nil {
		return err
	}
	if isAlreadyExists {
		return u.r.UpdateCartItemQuantity(c, cartItem)
	}
	return u.r.AddCartItem(c, cartItem)
}

func (u cartUsecaseImpl) DeleteFromCart(c context.Context, cartItem entity.CartItem) error {
	isExists, err := u.r.IsCartItemExistsByIDAndUserID(c, cartItem)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrDeleteFromCart, apperror.ErrCartItemNotExists, apperror.ErrCartItemNotExists)
	}
	err = u.r.DeleteCartItemByID(c, cartItem.ID)
	if err != nil {
		return err
	}
	return nil
}

func groupByPharmacy(items []entity.CartItem) []entity.GroupedCartItem {
	result := []entity.GroupedCartItem{}
	if len(items) == 0 {
		return result
	}

	group := &entity.GroupedCartItem{
		PharmacyID:   items[0].PharmacyID,
		PharmacyName: items[0].PharmacyName,
		Items:        []entity.CartItem{},
	}

	for _, item := range items {
		if group.PharmacyID != item.PharmacyID {
			result = append(result, *group)
			group = &entity.GroupedCartItem{
				PharmacyID:   item.PharmacyID,
				PharmacyName: item.PharmacyName,
				Items:        []entity.CartItem{},
			}
		}
		group.Items = append(group.Items, entity.CartItem{
			ID:                item.ID,
			UserID:            item.UserID,
			PharmacyProductID: item.PharmacyProductID,
			Name:              item.Name,
			Manufacturer:      item.Manufacturer,
			Image:             item.Image,
			Quantity:          item.Quantity,
			Subtotal:          item.Subtotal,
		})
	}
	
	return append(result, *group)
}

func (u cartUsecaseImpl) GetGroupedCartItems(c context.Context, userID int) ([]entity.GroupedCartItem, error) {
	cartItems, err := u.r.GetDetailedCartItemsByUserID(c, userID)
	if err != nil {
		return nil, err
	}
	return groupByPharmacy(cartItems), nil
}

func (u cartUsecaseImpl) GetCartItems(c context.Context, userID int) ([]entity.CartItem, error) {
	cartItems, err := u.r.GetCartItemsByUserID(c, userID)
	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (u cartUsecaseImpl) GetSelectedCartItems(c context.Context, userID int, ids []int) (*entity.ListGroupedCartItem, error) {
	cartItems, err := u.r.GetSelectedCartItems(c, userID, ids)
	if len(cartItems) == 0 {
		return nil, apperror.NewErrStatusBadRequest(appconstant.FieldErrGetCart, apperror.ErrCartNotAvailable, apperror.ErrCartNotAvailable)
	}
	if err != nil {
		return nil, err
	}
	for _, cartItem := range cartItems {
		stock, err := u.pp.GetStockByID(c, cartItem.PharmacyProductID)
		if err != nil {
			return nil, err
		}
		if cartItem.Quantity > stock {
			return nil, apperror.NewErrStatusBadRequest(appconstant.FieldErrGetCart, apperror.ErrStockUnavailable, apperror.ErrStockUnavailable)
		}
	}
	var listGrouped entity.ListGroupedCartItem
	selectedCarts := groupByPharmacy(cartItems)
	listGrouped.GroupedItem = selectedCarts
	listGrouped.ID = uuid.New().String()
	err = u.r.AddCheckoutItem(c, listGrouped, userID)
	if err != nil {
		return nil, err
	}
	return &listGrouped, nil
}
