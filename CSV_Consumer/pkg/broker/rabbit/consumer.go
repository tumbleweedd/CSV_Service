package rabbit

import (
	"encoding/csv"
	"fmt"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/model"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/repository"
	"log"
	"strings"
)

type Consumer struct {
	rabbitConn      *RabbitConnection
	storeRepository repository.Store
}

func NewConsumer(rabbitConn *RabbitConnection, repo repository.Store) *Consumer {
	return &Consumer{
		rabbitConn:      rabbitConn,
		storeRepository: repo,
	}
}

func (c *Consumer) ConsumeQueue(queueName string) error {
	messages, err := c.rabbitConn.ch.Consume(
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

	go func() {
		for message := range messages {

			user, err := convertMessageToUser(message.Body)
			if err != nil {
				log.Printf("Ошибка преобразования csv-like строки в структуру: %v", err)
			}

			if isAccepted, err := c.storeRepository.CheckForAccepted(user.Id); !isAccepted {
				if err := c.storeRepository.Save(user); err != nil {
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
