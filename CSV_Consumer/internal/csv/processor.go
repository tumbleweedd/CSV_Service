package csv

import (
	"github.com/sirupsen/logrus"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/repository"
	"github.com/tumbleweedd/intership/CSV_Consumer/pkg/broker/rabbit"
	"os"
	"strings"
)

func Process(fileSource string, rabbitCon *rabbit.RabbitConnection, done <-chan struct{}, repo *repository.Repository) {
	file, err := os.Open(fileSource)
	if err != nil {
		logrus.Infof("Ошибка os.Open: %s", err.Error())
		return
	}
	defer file.Close()

	consumer := rabbit.NewConsumer(rabbitCon, repo)
	defer rabbitCon.Close()
	queueName := strings.Split(fileSource, "\\")

	for {
		select {
		case <-done:
			return
		default:
			err := consumer.ConsumeQueue(queueName[2])
			if err != nil {
				logrus.Infof("Ошибка при прослушивании очереди: %s", err.Error())
			}
		}
	}

}
