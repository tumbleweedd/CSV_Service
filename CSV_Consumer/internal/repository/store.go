package repository

import (
	"context"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (sr *StoreRepository) Save(user *model.User) error {
	//const query = `insert into users (id, full_name, username, email, phone_number) VALUES ($1, $2, $3, $4, $5)`
	//
	//_, err := sr.db.Exec(query, user.Id, user.FullName, user.Username, user.Email, user.PhoneNumber)
	//
	//return err
	_, err := sr.db.InsertOne(context.TODO(), user)

	return err
}

func (sr *StoreRepository) CheckForAccepted(userId primitive.ObjectID) (bool, error) {
	//var isAccepted bool
	//
	//const query = `select u.accepted from users u where u.id = $1 `
	//
	//row := sr.db.QueryRow(query, userId)
	//
	//err := row.Scan(&isAccepted)
	//
	//return isAccepted, err

	return false, nil
}
