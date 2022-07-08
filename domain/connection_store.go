package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectionStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*ConnectionRequest, error)
	Insert(ctx context.Context, company *ConnectionRequest) error
	DeleteAll(ctx context.Context)
	GetRequestsForUser(ctx context.Context, username string) ([]*ConnectionRequest, error)
	Delete(ctx context.Context, id primitive.ObjectID)
	CheckIfUsersConnected(ctx context.Context, usernameFrom, usernameTo string) bool
	GetConnectionByUsernames(ctx context.Context, usernameFrom, usernameTo string) (*ConnectionRequest, error)
}
