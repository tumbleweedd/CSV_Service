package csv

import (
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/repository"
	"github.com/tumbleweedd/intership/CSV_Consumer/pkg/broker/rabbit"
	"github.com/tumbleweedd/intership/CSV_Consumer/pkg/logger"
	"os"
	"strings"
)

func Process(fileSource string, rabbitCon *rabbit.RabbitConnection, done <-chan struct{}, repo *repository.Repository, logger *logger.Logger) {
	file, err := os.Open(fileSource)
	if err != nil {
		logger.Infof("Ошибка os.Open: %s", err.Error())
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
				logger.Infof("Ошибка при прослушивании очереди: %s", err.Error())
			}
		}
	}

}
