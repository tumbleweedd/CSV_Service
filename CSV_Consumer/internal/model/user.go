package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Subscription struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" csv:"sub_id"`
	EventType string             `bson:"event_type" csv:"event_type"`
	Active    bool               `bson:"active" csv:"active"`
	CreatedAt time.Time          `bson:"created_at" csv:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" csv:"updated_at"`
}

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" csv:"user_id"`
	FullName         string             `bson:"full_name" csv:"full_name"`
	Username         string             `bson:"username" csv:"username"`
	Email            string             `bson:"email" csv:"email"`
	Phone            string             `bson:"phone" csv:"phone"`
	Telegram         string             `bson:"telegram" csv:"telegram"`
	Subscriptions    []Subscription     `bson:"subscriptions" csv:"subscriptions"`
	CreatedAt        time.Time          `bson:"created_at" csv:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" csv:"updated_at"`
	StatusOfDelivery bool               `bson:"status" csv:"status"`
}
