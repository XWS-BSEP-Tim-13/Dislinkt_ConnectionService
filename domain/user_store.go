package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	GetActiveById(ctx context.Context, id primitive.ObjectID) (*RegisteredUser, error)
	GetAllActive(ctx context.Context) ([]*RegisteredUser, error)
	GetActiveByUsername(ctx context.Context, username string) (*RegisteredUser, error)
	GetByUsername(ctx context.Context, username string) (*RegisteredUser, error)
	Update(ctx context.Context, user *RegisteredUser) error
	Insert(ctx context.Context, company *RegisteredUser) error
	DeleteAll(ctx context.Context)
	CheckIfUsersConnected(ctx context.Context, fromUsername, toUsername string) (*RegisteredUser, error)
	CheckIfUserIsBlocked(ctx context.Context, fromUsername, toUsername string) (*RegisteredUser, error)
	UpdateBlockedList(ctx context.Context, user *RegisteredUser) error
}
