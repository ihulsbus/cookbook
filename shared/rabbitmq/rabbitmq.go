package rabbitmq

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQ struct {
	Connection *rabbitmq.Conn
}

func NewRabbitMQConnection(username, password, host string, logger *log.Logger) (*RabbitMQ, error) {
	var connString string = fmt.Sprintf("amqp://%s:%s@%s", username, password, host)
	var conn *rabbitmq.Conn
	var err error

	conn, err = rabbitmq.NewConn(
		connString,
		rabbitmq.WithConnectionOptionsLogger(logger),
		rabbitmq.WithConnectionOptionsReconnectInterval(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{Connection: conn}, nil
}

func (r *RabbitMQ) Close() {
	if r.Connection != nil {
		r.Connection.Close()
	}
}
