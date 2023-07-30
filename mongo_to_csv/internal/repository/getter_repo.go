package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetterRepository struct {
	db *mongo.Collection
}

func NewGetterRepository(client *mongo.Client, dbName, collectionName string) *GetterRepository {
	db := client.Database(dbName).Collection(collectionName)
	return &GetterRepository{
		db: db,
	}
}

func (getterRepo *GetterRepository) GetUsersCursor(ctx context.Context) (*mongo.Cursor, error) {
	cur, err := getterRepo.db.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	return cur, nil
}

func (getterRepo *GetterRepository) GetUser(ctx context.Context, objectId primitive.ObjectID) *mongo.SingleResult {
	filter := bson.M{
		"_id": objectId,
	}

	res := getterRepo.db.FindOne(ctx, filter)

	return res
}

func (getterRepo *GetterRepository) GetSubsCursor(ctx context.Context, obj bson.D) (*mongo.Cursor, error) {
	cur, err := getterRepo.db.Find(ctx, obj)
	if err != nil {
		return nil, err
	}

	return cur, nil
}

func (getterRepo *GetterRepository) GetSub(ctx context.Context, objectId primitive.ObjectID) *mongo.SingleResult {
	filter := bson.M{
		"subscriptions._id": objectId,
	}
	res := getterRepo.db.FindOne(ctx, filter)

	return res
}
