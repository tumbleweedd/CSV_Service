package rabbitmq

import (
	"encoding/csv"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/model"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/repository/postgres"
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

func (r *RabbitConnection) ConsumeQueue(queueName string, repo postgres.Store) error {
	messages, err := r.ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	consumer := NewConsumer(r, repo)

	go func() {
		for message := range messages {

			user, err := convertMessageToUser(message.Body)
			if err != nil {
				log.Printf("Ошибка преобразования csv-like строки в структуру: %v", err)
			}

			if isAccepted, err := consumer.storeRepository.CheckForAccepted(user.Id); !isAccepted {
				if err := consumer.storeRepository.Save(user); err != nil {
					log.Printf("Ошибка при сохранении пользователя в базу: %v", err)
				}
				message.Ack(false)
			} else {
				switch err {
				case nil:
					continue
				default:
					log.Printf("Ошибка при проверке статуса доставки: %v", err)
				}
			}

		}
	}()

	return nil
}

func convertMessageToUser(message []byte) (*model.User, error) {
	strUser := string(message)

	reader := csv.NewReader(strings.NewReader(strUser))

	records, err := reader.Read()
	if err != nil {
		return nil, err
	}

	if len(records) != 5 {
		return nil, fmt.Errorf("неправильное количество полей в CSV строке")
	}

	user := &model.User{
		Id:          records[0],
		FullName:    records[1],
		Username:    records[2],
		Email:       records[3],
		PhoneNumber: records[4],
	}

	return user, err
}
