package repository

import (
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Save(user *model.User) error
	CheckForAccepted(userId primitive.ObjectID) (bool, error)
}

type Repository struct {
	Store
}

func NewStorage(mongoClient *mongo.Client, dbName, collectionName string) *Repository {
	return &Repository{
		Store: NewStoreRepository(mongoClient, dbName, collectionName),
	}
}
