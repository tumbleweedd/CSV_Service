package repository

import (
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"context"
)

type Store interface {
	Save(ctx context.Context, user *model.User) error
	FindOne(ctx context.Context, userId primitive.ObjectID, proj bson.M) (*model.User, error)
}

type Repository struct {
	Store
}

func NewStorage(mongoClient *mongo.Client, dbName, collectionName string) *Repository {
	return &Repository{
		Store: NewStoreRepository(mongoClient, dbName, collectionName),
	}
}
