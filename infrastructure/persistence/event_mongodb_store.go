package persistence

import (
	"context"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	EVENT_DATABASE   = "events"
	EVENT_COLLECTION = "event"
)

type EventMongoDBStore struct {
	events *mongo.Collection
}

func (store EventMongoDBStore) Get(id primitive.ObjectID) (*domain.Event, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store EventMongoDBStore) GetAll() ([]*domain.Event, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store EventMongoDBStore) Insert(event *domain.Event) error {
	result, err := store.events.InsertOne(context.TODO(), event)
	if err != nil {
		return err
	}
	event.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store EventMongoDBStore) DeleteAll() {
	store.events.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *EventMongoDBStore) filter(filter interface{}) ([]*domain.Event, error) {
	cursor, err := store.events.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decodeEvent(cursor)
}

func (store *EventMongoDBStore) filterOne(filter interface{}) (connection *domain.Event, err error) {
	result := store.events.FindOne(context.TODO(), filter)
	err = result.Decode(&connection)
	return
}

func decodeEvent(cursor *mongo.Cursor) (events []*domain.Event, err error) {
	for cursor.Next(context.TODO()) {
		var event domain.Event
		err = cursor.Decode(&event)
		if err != nil {
			return
		}
		events = append(events, &event)
	}
	err = cursor.Err()
	return
}

func NewEventMongoDBStore(client *mongo.Client) domain.EventStore {
	events := client.Database(EVENT_DATABASE).Collection(EVENT_COLLECTION)
	return &EventMongoDBStore{
		events: events,
	}
}
