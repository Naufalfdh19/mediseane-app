package usecase

import (
	"context"
	"encoding/json"
	"montelukast/pkg/logger"

	"github.com/streadway/amqp"
)

type RabbitMQConsumerStatus struct {
	usecase  OrderUsecase
	rabbitMQ *amqp.Channel
}

func NewRabbitMQConsumerStatus(rabbitMQ *amqp.Channel, u OrderUsecase) *RabbitMQConsumerStatus {
	return &RabbitMQConsumerStatus{usecase: u, rabbitMQ: rabbitMQ}
}

func (r *RabbitMQConsumerStatus) ConsumeDelayedMessage() {
	err := r.rabbitMQ.ExchangeDeclare(
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
		logger.Log.Error(err)

	}
	q, err := r.rabbitMQ.QueueDeclare(
		"update-status-queue",
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
		q.Name,                   // queue name
		"",                       // routing key
		"update-status-exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		logger.Log.Error(err)

	}
	msgs, err := r.rabbitMQ.Consume(
		q.Name,
		"status_update",
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
			ID int `json:"status_order"`
		}
		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			logger.Log.Error(err)
		}
		err = r.usecase.UpdateOrderStatusFromConsumer(context.Background(), data.ID)
		if err != nil {
			logger.Log.Error(err)
		}
		err = d.Ack(false)
		if err != nil {
			logger.Log.Error(err)
		}
	}
}
