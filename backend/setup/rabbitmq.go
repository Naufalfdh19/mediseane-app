package setup

import (
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"os"

	"github.com/streadway/amqp"
)

func ConnectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(os.Getenv("RABBIT_MQ_SERVER"))
	if err != nil {
		return nil, nil, apperror.NewErrInternalServerError(appconstant.FieldErrUploadImage, apperror.ErrInternalServer, err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, apperror.NewErrInternalServerError(appconstant.FieldErrUploadImage, apperror.ErrInternalServer, err)
	}
	return conn, ch, nil
}
