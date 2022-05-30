package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectionStore interface {
	Get(id primitive.ObjectID) (*Company, error)
	GetAll() ([]*Company, error)
	Insert(company *Company) error
	DeleteAll()
}
