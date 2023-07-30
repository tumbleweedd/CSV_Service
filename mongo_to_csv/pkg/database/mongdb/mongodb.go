package mongdb

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"fmt"
	"context"
)

type Config struct {
	MongoHost   string
	MongoPort   string
	MongoDBName string
}

type MongoSession struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoConnection(ctx context.Context, cfg *Config) (*MongoSession, error) {
	uri := fmt.Sprintf("mongodb://%s:%s", cfg.MongoHost, cfg.MongoPort)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	database := client.Database(cfg.MongoDBName)

	return &MongoSession{
		Client:   client,
		Database: database,
	}, nil
}

func (ms *MongoSession) disconnectSession(ctx context.Context) {
	defer func() {
		if err := ms.Client.Disconnect(ctx); err != nil {
			logrus.Fatalf("disconnec mongo session error:%s", err.Error())
		}
	}()
}
