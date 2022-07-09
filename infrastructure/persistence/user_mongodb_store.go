package persistence

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	USER_DATABASE   = "users"
	USER_COLLECTION = "user"
)

type UserMongoDBStore struct {
	users *mongo.Collection
}

func NewUserMongoDBStore(client *mongo.Client) domain.UserStore {
	users := client.Database(USER_DATABASE).Collection(USER_COLLECTION)
	return &UserMongoDBStore{
		users: users,
	}
}

func (store *UserMongoDBStore) GetActiveById(ctx context.Context, id primitive.ObjectID) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB GetActiveById")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"_id": id}
	return store.filterOne(ctx, filter)
}

func (store *UserMongoDBStore) GetAllActive(ctx context.Context) ([]*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB GetAllActive")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.D{{}}
	return store.filter(ctx, filter)
}

func (store *UserMongoDBStore) GetActiveByUsername(ctx context.Context, username string) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB GetActiveByUsername")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"username": username}
	return store.filterOne(ctx, filter)
}

func (store *UserMongoDBStore) GetByUsername(ctx context.Context, username string) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB GetByUsername")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"username": username}
	return store.filterOne(ctx, filter)
}

func (store *UserMongoDBStore) Update(ctx context.Context, user *domain.RegisteredUser) (err error) {
	span := tracer.StartSpanFromContext(ctx, "DB Update")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	fmt.Printf("Updating user %s %s\n", user.FirstName, user.Connections)
	filter := bson.M{"_id": user.Id}
	replacementObj := user
	_, err = store.users.ReplaceOne(ctx, filter, replacementObj)
	fmt.Printf("Updated \n")
	if err != nil {
		return err
	}
	return nil
}

func (store *UserMongoDBStore) Insert(ctx context.Context, user *domain.RegisteredUser) error {
	span := tracer.StartSpanFromContext(ctx, "DB Insert")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	result, err := store.users.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *UserMongoDBStore) CheckIfUsersConnected(ctx context.Context, fromUsername, toUsername string) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB CheckIfUsersConnected")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"connections": fromUsername, "username": toUsername}
	return store.filterOne(ctx, filter)
}

func (store *UserMongoDBStore) UpdateBlockedList(ctx context.Context, user *domain.RegisteredUser) error {
	span := tracer.StartSpanFromContext(ctx, "DB UpdateBlockedList")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	_, err := store.users.UpdateOne(
		ctx,
		bson.M{"_id": user.Id, "is_active": true},
		bson.D{
			{"$set", bson.D{{"blocked_users", user.BlockedUsers}}},
		},
	)
	return err
}

func (store *UserMongoDBStore) CheckIfUserIsBlocked(ctx context.Context, fromUsername, toUsername string) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB CheckIfUserIsBlocked")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"blocked_users": fromUsername, "username": toUsername}
	return store.filterOne(ctx, filter)
}

func (store *UserMongoDBStore) DeleteAll(ctx context.Context) {
	span := tracer.StartSpanFromContext(ctx, "DB DeleteAll")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	store.users.DeleteMany(ctx, bson.D{{}})
}

func (store *UserMongoDBStore) filter(ctx context.Context, filter interface{}) ([]*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB filter")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	cursor, err := store.users.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *UserMongoDBStore) filterOne(ctx context.Context, filter interface{}) (user *domain.RegisteredUser, err error) {
	span := tracer.StartSpanFromContext(ctx, "DB filterOne")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	result := store.users.FindOne(ctx, filter)
	err = result.Decode(&user)
	return
}

func decode(cursor *mongo.Cursor) (users []*domain.RegisteredUser, err error) {
	for cursor.Next(context.TODO()) {
		var user domain.RegisteredUser
		err = cursor.Decode(&user)
		if err != nil {
			return
		}
		users = append(users, &user)
	}
	err = cursor.Err()
	return
}
