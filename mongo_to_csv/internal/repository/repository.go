package repository

import "go.mongodb.org/mongo-driver/mongo"

type Getter interface {
}

type Repository struct {
	Getter
}

func NewRepository(mongoClient *mongo.Client, dbName, collectionName string) *Repository {
	return &Repository{
		Getter: NewGetterRepository(mongoClient, dbName, collectionName),
	}
}
