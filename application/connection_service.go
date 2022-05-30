package application

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectionService struct {
	store domain.ConnectionStore
}

func NewConnectionService(store domain.ConnectionStore) *ConnectionService {
	return &ConnectionService{
		store: store,
	}
}

func (service *ConnectionService) Get(id primitive.ObjectID) (*domain.Company, error) {
	return service.store.Get(id)
}

func (service *ConnectionService) GetAll() ([]*domain.Company, error) {
	return service.store.GetAll()
}
