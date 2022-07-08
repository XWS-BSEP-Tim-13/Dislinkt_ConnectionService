package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	Id        primitive.ObjectID `bson:"_id"`
	Action    string             `bson:"action"`
	User      string             `bson:"user"`
	Published time.Time          `bson:"published" validate:"required"`
}
