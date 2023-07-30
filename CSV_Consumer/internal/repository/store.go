package repository

import (
	"context"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StoreRepository struct {
	db *mongo.Collection
}

func NewStoreRepository(client *mongo.Client, dbName, collectionName string) *StoreRepository {
	db := client.Database(dbName).Collection(collectionName)
	return &StoreRepository{
		db: db,
	}
}

func (sr *StoreRepository) Save(ctx context.Context, user *model.User) error {
	_, err := sr.db.InsertOne(ctx, user)

	return err
}

func (sr *StoreRepository) FindOne(ctx context.Context, userId primitive.ObjectID, proj bson.M) (*model.User, error) {
	var user *model.User
	query := bson.M{"_id": userId}

	err := sr.db.FindOne(ctx, query, options.FindOne().SetProjection(proj)).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
