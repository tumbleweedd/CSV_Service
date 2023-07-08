package rabbitmq

import (
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/repository/postgres"
)

type Consumer struct {
	rabbitConn      *RabbitConnection
	storeRepository postgres.Store
}

func NewConsumer(rabbitConn *RabbitConnection, repo postgres.Store) *Consumer {
	return &Consumer{
		rabbitConn:      rabbitConn,
		storeRepository: repo,
	}
}
