package rabbit

import (
	"github.com/streadway/amqp"
)

type RabbitConnection struct {
	con *amqp.Connection
	ch  *amqp.Channel
}

func NewRabbitMQConnection(conURI string) (*RabbitConnection, error) {
	connectRabbitMQ, err := amqp.Dial(conURI)
	if err != nil {
		return nil, err
	}

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitConnection{
		con: connectRabbitMQ,
		ch:  channelRabbitMQ,
	}, nil
}

func (r *RabbitConnection) Close() {
	r.ch.Close()
	r.con.Close()
}
