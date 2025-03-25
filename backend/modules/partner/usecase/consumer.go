package usecase

import (
	"context"
	"encoding/json"
	"montelukast/modules/partner/entity"
	"montelukast/pkg/logger"

	"github.com/streadway/amqp"
)

type RabbitMQConsumerParter struct {
	usecase  PartnerUsecase
	rabbitMQ *amqp.Channel
}

func NewRabbitMQConsumerPartner(rabbitMQ *amqp.Channel, u PartnerUsecase) *RabbitMQConsumerParter {
	return &RabbitMQConsumerParter{usecase: u, rabbitMQ: rabbitMQ}
}

func (r *RabbitMQConsumerParter) ConsumeDelayedMessage() {

	err := r.rabbitMQ.ExchangeDeclare(
		"update-partner-exchange", //name
		"x-delayed-message",       //type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		amqp.Table{
			"x-delayed-type": "fanout",
		},
	)
	if err != nil {
		logger.Log.Error(err)

	}
	q, err := r.rabbitMQ.QueueDeclare(
		"update-partner-queue",
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
		q.Name,                    // queue name
		"",                        // routing key
		"update-partner-exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		logger.Log.Error(err)

	}
	msgs, err := r.rabbitMQ.Consume(
		q.Name,
		"parter_update",
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
		var wrapper struct {
			UpdatePartner entity.Partner `json:"update_partner"`
		}
		err := json.Unmarshal(d.Body, &wrapper)
		if err != nil {
			logger.Log.Error(err)
		}
		err = r.usecase.UpdatePartnerFromConsumer(context.Background(), wrapper.UpdatePartner)
		if err != nil {
			logger.Log.Error(err)
		}
		err = d.Ack(false)
		if err != nil {
			logger.Log.Error(err)
		}
	}
}
