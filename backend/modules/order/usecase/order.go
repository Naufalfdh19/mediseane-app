package usecase

import (
	"context"
	"encoding/json"
	"math"
	"montelukast/modules/order/entity"
	queryparams "montelukast/modules/order/query_params"
	"montelukast/modules/order/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"

	"github.com/streadway/amqp"
)

type OrderUsecase interface {
	GetOrders(c context.Context, queryParams queryparams.QueryParams, pharmacistID int) (*entity.OrdersList, error)
	GetOrderedProducts(c context.Context, orderDetailID int, pharmacistID int) (*entity.OrderProductDetail, error)
	DeleteOrder(c context.Context, orderDetailID int, pharmacistID int) error
	UpdateOrderStatus(c context.Context, orderDetailID int, pharmacistID int) error
	UpdateOrderStatusFromConsumer(c context.Context, orderDetailID int) error
}

type orderUsecaseImpl struct {
	r        repository.OrderRepo
	tr       transaction.TransactorRepoImpl
	rabbitMQ *amqp.Channel
}

func NewOrderUsecase(rabbitMQ *amqp.Channel, r repository.OrderRepo, tr transaction.TransactorRepoImpl) orderUsecaseImpl {
	return orderUsecaseImpl{
		r:        r,
		tr:       tr,
		rabbitMQ: rabbitMQ,
	}
}

func (u orderUsecaseImpl) GetOrders(c context.Context, queryParams queryparams.QueryParams, pharmacistID int) (*entity.OrdersList, error) {
	isExists, err := u.r.IsPharmacistExistsByID(c, pharmacistID)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, apperror.NewErrStatusNotFound(appconstant.FieldErrGetUserOrders, apperror.ErrUserNotExists, apperror.ErrUserNotExists)
	}

	pharmacyID, err := u.r.GetPharmacyIDByPharmacistID(c, pharmacistID)
	if err != nil {
		return nil, err
	}

	totalOrder, err := u.r.GetTotalOrder(c, queryParams, pharmacyID)
	if err != nil {
		return nil, err
	}

	if queryParams.Limit <= 0 {
		queryParams.Limit = 10
	}

	totalPage := int(math.Ceil(float64(totalOrder) / float64(queryParams.Limit)))
	if totalOrder <= 0 {
		totalPage = 1
	}

	queryParams.Page = apperror.CheckCurrentPage(queryParams.Page, totalPage, queryParams.Limit)

	orders, err := u.r.GetOrders(c, queryParams, pharmacyID)
	if err != nil {
		return nil, err
	}

	pagination := entity.Pagination{
		CurrentPage: queryParams.Page,
		TotalPage:   totalPage,
		TotalOrder:  totalOrder,
	}
	ordersList := entity.OrdersList{
		Pagination: pagination,
		Orders:     orders,
	}

	return &ordersList, nil
}

func (u orderUsecaseImpl) GetOrderedProducts(c context.Context, orderDetailID int, pharmacistID int) (*entity.OrderProductDetail, error) {
	var orderProductDetail entity.OrderProductDetail

	isExists, err := u.r.IsOrderDetailExistsByID(c, orderDetailID)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, apperror.NewErrStatusNotFound(appconstant.FieldErrGetOrderedProducts, apperror.ErrOrderDetailNotExists, apperror.ErrOrderDetailNotExists)
	}

	isExists, err = u.r.IsPharmacistExistsByID(c, pharmacistID)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, apperror.NewErrStatusNotFound(appconstant.FieldErrGetOrderedProducts, apperror.ErrPharmacistNotExists, apperror.ErrPharmacistNotExists)
	}

	isExists, err = u.r.IsPharmacyExistsByPharmacistID(c, pharmacistID)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, apperror.NewErrStatusNotFound(appconstant.FieldErrGetOrderedProducts, apperror.ErrPharmacyNotExists, apperror.ErrPharmacyNotExists)
	}

	isAuthorized, err := u.IsPharmacistAuthorized(c, orderDetailID, pharmacistID)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, apperror.NewErrStatusUnauthorized(appconstant.FieldErrGetOrderedProducts, apperror.ErrUserUnauthorized, apperror.ErrUserUnauthorized)
	}

	pharmacyID, err := u.r.GetPharmacyIDByOrderID(c, orderDetailID)
	if err != nil {
		return nil, err
	}

	orderDetail, err := u.r.GetOrderDetailByID(c, orderDetailID)
	if err != nil {
		return nil, err
	}

	productOrders, err := u.r.GetOrderedProduct(c, orderDetailID, pharmacyID)
	if err != nil {
		return nil, err
	}

	orderProductDetail.ID = orderDetail.ID
	orderProductDetail.Status = orderDetail.Status
	orderProductDetail.CreatedAt = orderDetail.CreatedAt
	orderProductDetail.ProductDetails = productOrders

	return &orderProductDetail, nil
}

func (u orderUsecaseImpl) DeleteOrder(c context.Context, orderDetailID int, pharmacistID int) error {
	err := u.tr.WithinTransaction(c, func(txCtx context.Context) error {
		isExists, err := u.r.IsOrderDetailExistsByID(txCtx, orderDetailID)
		if err != nil {
			return err
		}
		if !isExists {
			return apperror.NewErrStatusNotFound(appconstant.FieldErrDeleteOrder, apperror.ErrOrderDetailNotExists, apperror.ErrOrderDetailNotExists)
		}

		isExists, err = u.r.IsPharmacistExistsByID(txCtx, pharmacistID)
		if err != nil {
			return err
		}
		if !isExists {
			apperror.NewErrStatusNotFound(appconstant.FieldErrDeleteOrder, apperror.ErrOrderDetailNotExists, apperror.ErrOrderDetailNotExists)
		}

		isAuthorized, err := u.IsPharmacistAuthorized(c, orderDetailID, pharmacistID)
		if err != nil {
			return err
		}
		if !isAuthorized {
			return apperror.NewErrStatusUnauthorized(appconstant.FieldErrGetUserOrders, apperror.ErrUserUnauthorized, apperror.ErrUserUnauthorized)
		}

		orderStatus, err := u.r.GetOrderStatusByID(c, orderDetailID)
		if err != nil {
			return err
		}

		CanBeCanceled := apperror.IsOrderCanBeCanceled(orderStatus)
		if !CanBeCanceled {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrDeleteOrder, apperror.ErrOrderCannotCanceled, apperror.ErrOrderCannotCanceled)
		}

		err = u.r.DeleteOrderDetails(txCtx, orderDetailID)
		if err != nil {
			return err
		}

		err = u.r.DeleteOrderProductDetails(txCtx, orderDetailID)
		if err != nil {
			return err
		}

		orderedProducts, err := u.r.GetOrderedProductsQuantity(txCtx, orderDetailID)
		if err != nil {
			return err
		}

		for pharmacyProductID, quantity := range orderedProducts {
			err := u.r.UpdateProductQuantity(txCtx, pharmacyProductID, quantity)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u orderUsecaseImpl) UpdateOrderStatus(c context.Context, orderDetailID int, pharmacistID int) error {

	isExists, err := u.r.IsOrderDetailExistsByID(c, orderDetailID)
	if err != nil {
		return err
	}
	if !isExists {
		return apperror.NewErrStatusNotFound(appconstant.FieldErrDeleteOrder, apperror.ErrOrderDetailNotExists, apperror.ErrOrderDetailNotExists)
	}

	isExists, err = u.r.IsPharmacistExistsByID(c, pharmacistID)
	if err != nil {
		return err
	}
	if !isExists {
		apperror.NewErrStatusNotFound(appconstant.FieldErrDeleteOrder, apperror.ErrOrderDetailNotExists, apperror.ErrOrderDetailNotExists)
	}

	isAuthorized, err := u.IsPharmacistAuthorized(c, orderDetailID, pharmacistID)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return apperror.NewErrStatusUnauthorized(appconstant.FieldErrGetUserOrders, apperror.ErrUserUnauthorized, apperror.ErrUserUnauthorized)
	}

	orderStatus, err := u.r.GetOrderStatusByID(c, orderDetailID)
	if err != nil {
		return err
	}

	if orderStatus != appconstant.StatusProcessing {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrUpdateOrderStatus, apperror.ErrCannotUpdateOrderStatus, apperror.ErrCannotUpdateOrderStatus)
	}
	err = u.r.UpdateOrderStatus(c, orderDetailID)
	if err != nil {
		return err
	}
	err = u.PublishDelayedMessage(c, orderDetailID, appconstant.OrderCompleteTime)
	if err != nil {
		return err
	}
	return nil
}

func (u orderUsecaseImpl) UpdateOrderStatusFromConsumer(c context.Context, orderDetailID int) error {
	err := u.r.UpdateOrderStatusDelivered(c, orderDetailID)
	if err != nil {
		return err
	}
	return nil
}

func (u orderUsecaseImpl) PublishDelayedMessage(c context.Context, id int, delay int) error {
	err := u.rabbitMQ.ExchangeDeclare(
		"update-status-exchange", //name
		"x-delayed-message",      //type
		true,                     // durable
		false,                    // auto-deleted
		false,                    // internal
		false,                    // no-wait
		amqp.Table{
			"x-delayed-type": "fanout",
		},
	)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrChangeStatus, apperror.ErrInternalServer, err)
	}
	body, err := json.Marshal(map[string]interface{}{
		"status_order": id,
	})

	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrChangeStatus, apperror.ErrInternalServer, err)
	}
	err = u.rabbitMQ.Publish(
		"update-status-exchange",
		"x-delayed-message",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Headers: amqp.Table{
				"x-delay": delay,
			},
		},
	)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrChangeStatus, apperror.ErrInternalServer, err)
	}
	return nil
}

func (u orderUsecaseImpl) IsPharmacistAuthorized(c context.Context, orderDetailID int, pharmacistID int) (bool, error) {
	pharmacyID1, err := u.r.GetPharmacyIDByPharmacistID(c, pharmacistID)
	if err != nil {
		return false, err
	}
	pharmacistID2, err := u.r.GetPharmacyIDByOrderID(c, orderDetailID)
	if err != nil {
		return false, err
	}
	return pharmacyID1 == pharmacistID2, nil
}
