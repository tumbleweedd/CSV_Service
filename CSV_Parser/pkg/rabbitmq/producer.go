package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"strings"
)

type Producer struct {
	rabbit *RabbitConnection
}

func NewProducer(rabbit *RabbitConnection) *Producer {
	return &Producer{
		rabbit: rabbit,
	}
}

func (p *Producer) SendToQueue(queueName string, data []string) error {
	_, err := p.rabbit.ch.QueueDeclare(
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

	if err = p.rabbit.ch.Publish(
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
