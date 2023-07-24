package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Subscription struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	EventType string             `bson:"eventType"`
	Active    bool               `bson:"active"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type User struct {
	ID            primitive.ObjectID `csv:"-" bson:"_id,omitempty"`
	FullName      string             `csv:"full_name" bson:"full_name"`
	Username      string             `bson:"username"`
	Email         string             `bson:"email"`
	Phone         string             `bson:"phone"`
	Telegram      string             `bson:"telegram"`
	Subscriptions []Subscription     `bson:"subscriptions"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}
