package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"strings"
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

func (r *RabbitConnection) SendToQueue(queueName string, data []string) error {
	_, err := r.ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(strings.Join(data, ",")),
	}

	if err = r.ch.Publish(
		"",
		queueName,
		false,
		false,
		message,
	); err != nil {
		return err
	}

	log.Println("Данные успешно отправлены в RabbitMQ:", string(message.Body))
	return nil
}
