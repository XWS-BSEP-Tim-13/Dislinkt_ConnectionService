package application

import (
	"errors"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventService struct {
	store domain.EventStore
}

func NewEventService(store domain.EventStore) *EventService {
	return &EventService{}
}

func (service *EventService) Get(id primitive.ObjectID) (*domain.Event, error) {
	return service.store.Get(id)
}

func (service *EventService) GetAll() ([]*domain.Event, error) {
	return service.store.GetAll()
}

func (service *EventService) CreateNewEvent(event *domain.Event) (*domain.Event, error) {
	(*event).Id = primitive.NewObjectID()
	err := service.store.Insert(event)
	if err != nil {
		err := errors.New("error while creating new event")
		return nil, err
	}

	return event, nil
}
