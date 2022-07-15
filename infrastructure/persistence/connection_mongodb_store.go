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
	DATABASE   = "connections"
	COLLECTION = "connection"
)

type ConnectionsMongoDBStore struct {
	connections *mongo.Collection
}

func NewConnectionMongoDBStore(client *mongo.Client) *ConnectionsMongoDBStore {
	companies := client.Database(DATABASE).Collection(COLLECTION)
	return &ConnectionsMongoDBStore{
		connections: companies,
	}
}

func (store ConnectionsMongoDBStore) Delete(ctx context.Context, id primitive.ObjectID) {
	span := tracer.StartSpanFromContext(ctx, "DB Delete")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"_id": id}
	_, err := store.connections.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
}

func (store ConnectionsMongoDBStore) GetRequestsForUser(ctx context.Context, username string) ([]*domain.ConnectionRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "DB GetRequestsForUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.D{{"to.username", username}}
	return store.filter(ctx, filter)
}

func (store ConnectionsMongoDBStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.ConnectionRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"_id": id}
	return store.filterOne(ctx, filter)
}

func (store ConnectionsMongoDBStore) Insert(ctx context.Context, connection *domain.ConnectionRequest) error {
	span := tracer.StartSpanFromContext(ctx, "DB Insert")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	result, err := store.connections.InsertOne(ctx, connection)
	if err != nil {
		return err
	}
	connection.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store ConnectionsMongoDBStore) DeleteAll(ctx context.Context) {
	span := tracer.StartSpanFromContext(ctx, "DB DeleteAll")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	store.connections.DeleteMany(ctx, bson.D{{}})
}

func (store ConnectionsMongoDBStore) GetConnectionByUsernames(ctx context.Context, usernameFrom, usernameTo string) (*domain.ConnectionRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "DB GetConnectionByUsernames")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"to.username": usernameTo, "from.username": usernameFrom}
	return store.filterOne(ctx, filter)
}

func (store ConnectionsMongoDBStore) CheckIfUsersConnected(ctx context.Context, usernameFrom, usernameTo string) bool {
	span := tracer.StartSpanFromContext(ctx, "DB CheckIfUsersConnected")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"to.username": usernameTo, "from.username": usernameFrom}
	_, err := store.filterOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (store *ConnectionsMongoDBStore) filter(ctx context.Context, filter interface{}) ([]*domain.ConnectionRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "DB filter")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	cursor, err := store.connections.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}
	return decodeConnection(cursor)
}

func (store *ConnectionsMongoDBStore) filterOne(ctx context.Context, filter interface{}) (connection *domain.ConnectionRequest, err error) {
	span := tracer.StartSpanFromContext(ctx, "DB filterOne")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	result := store.connections.FindOne(ctx, filter)
	err = result.Decode(&connection)
	return
}

func decodeConnection(cursor *mongo.Cursor) (connections []*domain.ConnectionRequest, err error) {
	for cursor.Next(context.TODO()) {
		var connection domain.ConnectionRequest
		err = cursor.Decode(&connection)
		if err != nil {
			return
		}
		connections = append(connections, &connection)
	}
	err = cursor.Err()
	return
}
