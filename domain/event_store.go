package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventStore interface {
	Get(id primitive.ObjectID) (*Event, error)
	GetAll() ([]*Event, error)
	Insert(event *Event) error
	DeleteAll()
}
