package usecase

import (
	"context"
	"fmt"
	"montelukast/modules/checkout/entity"
	"montelukast/modules/checkout/repository"
	delivery "montelukast/modules/delivery/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"

	"github.com/shopspring/decimal"
)

type CheckoutUsecase interface {
	Checkout(c context.Context, checkoutData entity.CheckoutData, userID int) <-chan error
	CancelOrderByUser(c context.Context, userID int, orderID int) error
}

type checkoutUsecaseImpl struct {
	c  repository.CheckoutRepo
	d  delivery.DeliveryRepository
	tr transaction.TransactorRepoImpl
}

func NewCheckoutUsecase(c repository.CheckoutRepo, d delivery.DeliveryRepository, tr transaction.TransactorRepoImpl) CheckoutUsecase {
	return checkoutUsecaseImpl{
		tr: tr,
		c:  c,
		d:  d,
	}
}

func (u checkoutUsecaseImpl) Checkout(c context.Context, checkoutData entity.CheckoutData, userID int) <-chan error {
	output := make(chan error)
	go func() {
		defer close(output)
		_, addressID, err := u.d.GetUserPostalCode(c, userID)
		if err != nil {
			output <- err
			return
		}
		result, err := u.c.GetCheckoutCartRedis(c, checkoutData.IDCart, userID)
		if err != nil {
			output <- err
			return
		}
		deliveryDict := make(map[int]int)
		if len(result.GroupedItem) != len(checkoutData.ListDeliveryData) {
			output <- apperror.NewErrStatusBadRequest(appconstant.FieldErrCheckout, apperror.ErrInvalidDeliveryData, apperror.ErrInvalidDeliveryData)
			return
		}
		for _, delivery := range checkoutData.ListDeliveryData {
			if _, ok := deliveryDict[delivery.PharmacyID]; !ok {
				deliveryDict[delivery.PharmacyID] = delivery.DeliveryID
			}
		}
		var listPrice []entity.DeliveryPriceData
		for _, pharmacy := range result.GroupedItem {
			list_ongkir, err := u.d.GetListOngkir(c, addressID, pharmacy.PharmacyID)
			if err != nil {
				output <- err
				return
			}
			for _, ongkir := range list_ongkir {
				data := entity.DeliveryPriceData{
					PharmacyID:    pharmacy.PharmacyID,
					LogisticPrice: ongkir.Cost,
					Status:        appconstant.StatusPending,
				}
				if deliveryDict[pharmacy.PharmacyID] == ongkir.Id {
					listPrice = append(listPrice, data)
				}
			}
		}
		var totalPrice decimal.Decimal
		for _, cart := range result.GroupedItem {
			for _, cost := range listPrice {
				if cost.PharmacyID == cart.PharmacyID {
					totalPrice = totalPrice.Add(cost.LogisticPrice)
				}
			}
			for _, item := range cart.Items {
				totalPrice = totalPrice.Add(item.Subtotal)
			}
		}
		err = u.tr.WithinTransaction(c, func(txCtx context.Context) error {
			idOrder, err := u.c.AddOrder(txCtx, totalPrice, userID)
			if err != nil {
				return err
			}
			ids, err := u.c.AddCheckoutOrderDetail(txCtx, listPrice, idOrder)
			if err != nil {
				return err
			}
			var productUnavailable []string
			var productInactive []string
			for i, pharmacy := range result.GroupedItem {
				for _, product := range pharmacy.Items {
					isExists, err := u.c.IsCartItemExistsByIDAndUserID(c, product)
					if err != nil {
						return err
					}
					if !isExists {
						return apperror.NewErrStatusNotFound(appconstant.FieldErrDeleteFromCart, apperror.ErrCartItemNotExists, apperror.ErrCartItemNotExists)
					}
					err = u.c.DeleteCartItemByID(txCtx, product.ID)
					if err != nil {
						return err
					}
					err = u.c.UpdateStock(txCtx, product)
					if err != nil {
						return err
					}
					quantity, err := u.c.GetQuantity(txCtx, product.PharmacyProductID)
					if err != nil {
						return err
					}
					is_active, err := u.c.IsActive(txCtx, product.PharmacyProductID)
					if err != nil {
						return err
					}
					if !is_active {
						productInactive = append(productInactive, product.Name)
					}
					if quantity < 0 {
						productUnavailable = append(productUnavailable, product.Name)
					}
				}
				err := u.c.AddOrderProductDetails(txCtx, ids[i], pharmacy)
				if err != nil {
					return err
				}
			}
			if len(productInactive) > 0 {
				return apperror.NewErrStatusBadRequest(fmt.Sprintf("Items not available %s", productInactive), apperror.ErrStockUnavailable, apperror.ErrStockUnavailable)
			}
			if len(productUnavailable) > 0 {
				return apperror.NewErrStatusBadRequest(fmt.Sprintf("Items stock available %s", productUnavailable), apperror.ErrStockUnavailable, apperror.ErrStockUnavailable)
			}
			err = u.c.DeleteCartRedis(c, checkoutData.IDCart, userID)
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			output <- err
			return
		}
		output <- nil
	}()
	return output
}

func (u checkoutUsecaseImpl) CancelOrderByUser(c context.Context, userID int, orderID int) error {
	exists, err := u.c.IsOrderExistByID(c, orderID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrCancel, apperror.ErrInvalidOrderCancelation, apperror.ErrInvalidOrderCancelation)
	}
	isOrderFromUser, err := u.c.IsOrderFromUser(c, orderID, userID)
	if err != nil {
		return err
	}
	if !isOrderFromUser {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrCancel, apperror.ErrInvalidOrderCancelation, apperror.ErrInvalidOrderCancelation)
	}
	err = u.c.CancelOrder(c, orderID)
	if err != nil {
		return err
	}
	return nil
}
