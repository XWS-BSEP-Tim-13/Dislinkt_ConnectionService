package persistence

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
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

func (store *UserMongoDBStore) GetActiveById(id primitive.ObjectID) (*domain.RegisteredUser, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetAllActive() ([]*domain.RegisteredUser, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *UserMongoDBStore) GetActiveByUsername(username string) (*domain.RegisteredUser, error) {
	filter := bson.M{"username": username}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetByUsername(username string) (*domain.RegisteredUser, error) {
	filter := bson.M{"username": username}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) Update(user *domain.RegisteredUser) (err error) {
	fmt.Printf("Updating user %s %s\n", user.FirstName, user.Connections)
	filter := bson.M{"_id": user.Id}
	replacementObj := user
	_, err = store.users.ReplaceOne(context.TODO(), filter, replacementObj)
	fmt.Printf("Updated \n")
	if err != nil {
		return err
	}
	return nil
}

func (store *UserMongoDBStore) Insert(user *domain.RegisteredUser) error {
	result, err := store.users.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	user.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *UserMongoDBStore) DeleteAll() {
	store.users.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *UserMongoDBStore) filter(filter interface{}) ([]*domain.RegisteredUser, error) {
	cursor, err := store.users.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *UserMongoDBStore) filterOne(filter interface{}) (user *domain.RegisteredUser, err error) {
	result := store.users.FindOne(context.TODO(), filter)
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
