package rabbit

import (
	"encoding/csv"
	"fmt"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/model"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/repository"
	"strings"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"strconv"
	"go.mongodb.org/mongo-driver/bson"
	"context"
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
			var usersID []primitive.ObjectID
			user, err := convertMessageToUser(message.Body, &usersID)
			if err != nil {
				logrus.Printf("Ошибка преобразования csv-like строки в структуру: %v", err)
			}

			proj := bson.M{"status": 1}
			if userForStatusCheck, err := c.storeRepository.FindOne(context.TODO(), user.ID, proj); !userForStatusCheck.StatusOfDelivery {
				if err := c.storeRepository.Save(context.TODO(), user); err != nil {
					logrus.Printf("Ошибка при сохранении пользователя в базу: %v", err)
				}
				message.Ack(false)
			} else {
				switch err {
				case nil:
					continue
				default:
					logrus.Printf("Ошибка при проверке статуса доставки: %v", err)
				}
			}

		}
	}()

	return nil
}

func convertMessageToUser(message []byte, usersID *[]primitive.ObjectID) (*model.User, error) {
	strUser := string(message)
	reader := csv.NewReader(strings.NewReader(strUser))
	records, err := reader.Read()
	if err != nil {
		return nil, err
	}

	if len(records) != 10 {
		return nil, fmt.Errorf("неправильное количество полей в CSV строке")
	}

	userId, err := primitive.ObjectIDFromHex(records[0])
	if err != nil {
		return nil, fmt.Errorf("failed parse userID %q to ObjectID:%v", records[0], err)
	}

	var currentUser *model.User
	for _, id := range *usersID {
		if id == userId {
			currentUser.ID = id
			break
		}
	}

	if currentUser == nil {
		createdAt, err := time.Parse(time.RFC3339, records[7])
		if err != nil {
			return nil, err
		}

		updatedAt, err := time.Parse(time.RFC3339, records[8])
		if err != nil {
			return nil, err
		}
		currentUser = &model.User{
			ID:        userId,
			FullName:  records[1],
			Username:  records[2],
			Email:     records[3],
			Phone:     records[4],
			Telegram:  records[5],
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		*usersID = append(*usersID, userId)
	}

	subscription, err := parseSubscription(records)
	if err != nil {
		return nil, err
	}

	currentUser.Subscriptions = append(currentUser.Subscriptions, *subscription)
	return currentUser, err
}

func parseSubscription(records []string) (*model.Subscription, error) {
	var subscription *model.Subscription

	subscriptionsStr := strings.Split(records[6], ";")
	for _, subStr := range subscriptionsStr {
		subFields := strings.Split(subStr, ",")
		if len(subFields) < 5 {
			return nil, fmt.Errorf("неправильное количество полей для подписки в CSV строке")
		}
		subID, err := primitive.ObjectIDFromHex(subFields[0])
		if err != nil {
			return nil, err
		}
		eventType := subFields[1]
		isActive, err := strconv.ParseBool(subFields[2])
		if err != nil {
			return nil, err
		}
		createdAt, err := time.Parse(time.RFC3339, subFields[3])
		if err != nil {
			return nil, err
		}
		updatedAt, err := time.Parse(time.RFC3339, subFields[4])
		if err != nil {
			return nil, err
		}

		subscription = &model.Subscription{
			ID:        subID,
			EventType: eventType,
			Active:    isActive,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return subscription, nil
}
