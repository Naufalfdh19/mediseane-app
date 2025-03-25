package usecase

import (
	"context"
	"encoding/json"
	userRepo "montelukast/modules/user/repository"
	"montelukast/modules/userorder/entity"
	"montelukast/modules/userorder/repository"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/transaction"

	"github.com/go-playground/validator/v10"
	"github.com/streadway/amqp"
)

type UserOrderUsecase interface {
	UpdatePaymentStatus(c context.Context, file entity.File, orderID int, userID int) error
	UpdateDeliveryStatus(c context.Context, order entity.Order) error
	GetDetailedOrders(c context.Context, filter entity.OrderFilter) ([]entity.Order, error)
	UpdatePaymentStatusFromConsumer(ctx context.Context, orderDetailIDs []int) error
}

type userOrderUsecaseImpl struct {
	r        repository.UserOrderRepo
	tr       transaction.TransactorRepoImpl
	userRepo userRepo.UserRepoImpl
	rabbitMQ *amqp.Channel
}

func NewUserOrderUsecase(rabbitMQ *amqp.Channel, r repository.UserOrderRepo, tr transaction.TransactorRepoImpl, userRepo userRepo.UserRepoImpl) userOrderUsecaseImpl {
	return userOrderUsecaseImpl{
		r:        r,
		tr:       tr,
		userRepo: userRepo,
		rabbitMQ: rabbitMQ,
	}
}

func (u userOrderUsecaseImpl) UpdatePaymentStatus(c context.Context, file entity.File, orderID int, userID int) error {
	validate := validator.New()
	err := validate.Struct(file)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrUploadImage, apperror.ErrInternalServer, err)
	}
	exists, err := u.r.IsOrderExistsByID(c, orderID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrUploadPayment, apperror.ErrOrderNotExists, apperror.ErrOrderNotExists)
	}

	orderDetailIDs, err := u.r.GetOrderDetailIDByID(c, orderID)
	if err != nil {
		return err
	}
	if len(orderDetailIDs) == 0 {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrUploadPayment, apperror.ErrOrderNotExists, apperror.ErrOrderNotExists)
	}

	allStatus, err := u.r.GetStatusesByID(c, orderID)
	if err != nil {
		return err
	}
	for _, status := range allStatus {
		if status != appconstant.StatusPending {
			return apperror.NewErrStatusBadRequest(appconstant.FieldErrUploadPayment, apperror.ErrPaymentAlreadyDone, apperror.ErrPaymentAlreadyDone)
		}
	}
	err = u.PublishDelayedMessage(c, orderDetailIDs, appconstant.PaymentCompleteTime)
	if err != nil {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrUploadPayment, apperror.ErrInternalServer, err)
	}
	return nil
}

func (u userOrderUsecaseImpl) PublishDelayedMessage(c context.Context, orderDetailIDs []int, delay int) error {
	err := u.rabbitMQ.ExchangeDeclare(
		"update-userorder-exchange", //name
		"x-delayed-message",         //type
		true,                        // durable
		false,                       // auto-deleted
		false,                       // internal
		false,                       // no-wait
		amqp.Table{
			"x-delayed-type": "fanout",
		},
	)
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrUpdatePartner, apperror.ErrInternalServer, err)
	}
	body, err := json.Marshal(map[string]interface{}{
		"order_detail_ids": orderDetailIDs,
	})
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrUploadImage, apperror.ErrInternalServer, err)
	}
	err = u.rabbitMQ.Publish(
		"update-userorder-exchange",
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
		return apperror.NewErrInternalServerError(appconstant.FieldErrUploadImage, apperror.ErrInternalServer, err)
	}
	return nil
}

func (u userOrderUsecaseImpl) UpdatePaymentStatusFromConsumer(ctx context.Context, orderDetailIDs []int) error {
	err := u.tr.WithinTransaction(ctx, func(txCtx context.Context) error {
		for _, orderDetailID := range orderDetailIDs {
			err := u.r.UpdatePaymentStatus(txCtx, orderDetailID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return apperror.NewErrInternalServerError(appconstant.FieldErrUploadImage, apperror.ErrInternalServer, err)
	}
	return nil
}

func (u userOrderUsecaseImpl) UpdateDeliveryStatus(c context.Context, order entity.Order) error {
	exists, err := u.r.IsOrderDetailExists(c, order)
	if err != nil {
		return err
	}
	if !exists {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrConfirmDelivery, apperror.ErrOrderNotExists, apperror.ErrOrderNotExists)
	}

	status, err := u.r.GetStatusByOrderDetailID(c, order.OrderDetails[0].ID)
	if err != nil {
		return err
	}
	if status != appconstant.StatusShipped {
		return apperror.NewErrStatusBadRequest(appconstant.FieldErrConfirmDelivery, apperror.ErrOrderCannotBeCompleted, apperror.ErrOrderCannotBeCompleted)
	}

	err = u.r.UpdateDeliveryStatusByOrderDetailID(c, order.OrderDetails[0].ID)
	if err != nil {
		return err
	}

	return nil
}

func groupByOrderDetail(rows []entity.OrderProductDetail) []entity.OrderDetail {
	grouped := make(map[int]*entity.OrderDetail)
	for _, row := range rows {
		if _, exists := grouped[row.OrderDetail.ID]; !exists {
			order := &entity.Order{
				ID:         row.OrderDetail.Order.ID,
				UserID:     row.OrderDetail.Order.UserID,
				TotalPrice: row.OrderDetail.Order.TotalPrice,
				CreatedAt:  row.OrderDetail.Order.CreatedAt,
			}
			grouped[row.OrderDetail.ID] = &entity.OrderDetail{
				ID:                  row.OrderDetail.ID,
				OrderID:             row.OrderDetail.OrderID,
				PharmacyID:          row.OrderDetail.PharmacyID,
				PharmacyName:        row.OrderDetail.PharmacyName,
				Status:              row.OrderDetail.Status,
				LogisticPrice:       row.OrderDetail.LogisticPrice,
				OrderProductDetails: []entity.OrderProductDetail{},
				Order:               *order,
			}
		}
		grouped[row.OrderDetail.ID].OrderProductDetails = append(grouped[row.OrderDetail.ID].OrderProductDetails, entity.OrderProductDetail{
			ID:                row.ID,
			PharmacyProductID: row.PharmacyProductID,
			Name:              row.Name,
			Manufacturer:      row.Manufacturer,
			Image:             row.Image,
			Quantity:          row.Quantity,
			Subtotal:          row.Subtotal,
		})
	}
	result := []entity.OrderDetail{}
	for _, group := range grouped {
		result = append(result, *group)
	}
	return result
}

func groupByOrder(rows []entity.OrderDetail) []entity.Order {
	grouped := make(map[int]*entity.Order)
	for _, row := range rows {
		if _, exists := grouped[row.Order.ID]; !exists {
			grouped[row.Order.ID] = &entity.Order{
				ID:           row.Order.ID,
				UserID:       row.Order.UserID,
				TotalPrice:   row.Order.TotalPrice,
				CreatedAt:    row.Order.CreatedAt,
				OrderDetails: []entity.OrderDetail{},
			}
		}
		grouped[row.Order.ID].OrderDetails = append(grouped[row.Order.ID].OrderDetails, row)
	}
	result := []entity.Order{}
	for _, group := range grouped {
		result = append(result, *group)
	}
	return result
}

func (u userOrderUsecaseImpl) GetDetailedOrders(c context.Context, filter entity.OrderFilter) ([]entity.Order, error) {
	orderProductDetails, err := u.r.GetDetailedOrdersByUserID(c, filter)
	if err != nil {
		return nil, err
	}
	groupedByOrderDetail := groupByOrderDetail(orderProductDetails)
	groupedByOrder := groupByOrder(groupedByOrderDetail)
	return groupedByOrder, nil
}
