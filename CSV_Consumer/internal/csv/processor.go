package csv

import (
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/rabbitmq"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/repository/postgres"
	"github.com/tumbleweedd/intership/CSV_Consumer/pkg/logger"
	"os"
	"strings"
)

func Process(fileSource string, rabbitMQDSN string, done <-chan struct{}, repo *postgres.Repository, logger *logger.Logger) {
	file, err := os.Open(fileSource)
	if err != nil {
		logger.Infof("Ошибка os.Open: %s", err.Error())
		return
	}
	defer file.Close()

	rabbit, err := rabbitmq.NewRabbitMQConnection(rabbitMQDSN)
	if err != nil {
		logger.Infof("Ошибка при инициализации RabbitConnection: %s", err.Error())
		return
	}
	defer rabbit.Close()
	queueName := strings.Split(fileSource, "\\")

	for {
		select {
		case <-done:
			return
		default:
			err := rabbit.ConsumeQueue(queueName[2], repo)
			if err != nil {
				logger.Infof("Ошибка при прослушивании очереди: %s", err.Error())
			}
		}
	}

}
