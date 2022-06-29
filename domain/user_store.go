package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserStore interface {
	GetActiveById(id primitive.ObjectID) (*RegisteredUser, error)
	GetAllActive() ([]*RegisteredUser, error)
	GetActiveByUsername(username string) (*RegisteredUser, error)
	GetByUsername(username string) (*RegisteredUser, error)
	Update(user *RegisteredUser) error
	Insert(company *RegisteredUser) error
	DeleteAll()
	CheckIfUsersConnected(fromUsername, toUsername string) (*RegisteredUser, error)
	CheckIfUserIsBlocked(fromUsername, toUsername string) (*RegisteredUser, error)
	UpdateBlockedList(user *RegisteredUser) error
}
