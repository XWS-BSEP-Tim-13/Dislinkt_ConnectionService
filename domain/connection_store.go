package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectionStore interface {
	Get(id primitive.ObjectID) (*ConnectionRequest, error)
	Insert(company *ConnectionRequest) error
	DeleteAll()
	GetRequestsForUser(username string) ([]*ConnectionRequest, error)
	Delete(id primitive.ObjectID)
}
