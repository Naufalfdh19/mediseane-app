package usecase

import (
	"context"
	"encoding/json"
	"montelukast/pkg/logger"

	"github.com/streadway/amqp"
)

type RabbitMQConsumer struct {
	usecase  UserOrderUsecase
	rabbitMQ *amqp.Channel
}

func NewRabbitMQConsumer(rabbitMQ *amqp.Channel, uc UserOrderUsecase) *RabbitMQConsumer {
	return &RabbitMQConsumer{usecase: uc, rabbitMQ: rabbitMQ}
}

func (r *RabbitMQConsumer) ConsumeDelayedMessage() {
	err := r.rabbitMQ.ExchangeDeclare(
		"update-userorder-exchange", //name
		"x-delayed-message",         //type
		true,                        // durable
		false,                       // auto-deleted
		false,                       // internal
		false,                       // no-wait
		amqp.Table{
			"x-delayed-type": "fanout",
		}, //arguments
	)
	if err != nil {
		logger.Log.Error(err)
	}
	q, err := r.rabbitMQ.QueueDeclare(
		"update-userorder-queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Log.Error(err)
	}
	err = r.rabbitMQ.QueueBind(
		q.Name,                      // queue name
		"",                          // routing key
		"update-userorder-exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		logger.Log.Error(err)
	}
	msgs, err := r.rabbitMQ.Consume(
		q.Name,
		"payment_update",
		false, //auto ack
		false, //exclusive
		false, //no-local
		false, //no-wait
		nil,
	)
	if err != nil {
		logger.Log.Error(err)

	}
	for d := range msgs {
		var data struct {
			OrderDetailIDs []int `json:"order_detail_ids"`
		}
		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			logger.Log.Error(err)
		}
		err = r.usecase.UpdatePaymentStatusFromConsumer(context.Background(), data.OrderDetailIDs)
		if err != nil {
			logger.Log.Error(err)

		}
		err = d.Ack(false)
		if err != nil {
			logger.Log.Error(err)

		}
	}
}
