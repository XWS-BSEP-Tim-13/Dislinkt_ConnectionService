package application

import (
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ConnectionService struct {
	store           domain.ConnectionStore
	userStore       domain.UserStore
	connectionNeo4j persistence.ConnectionNeo4jStore
}

func NewConnectionService(store domain.ConnectionStore, userStore domain.UserStore, neo4j persistence.ConnectionNeo4jStore) *ConnectionService {
	return &ConnectionService{
		store:           store,
		userStore:       userStore,
		connectionNeo4j: neo4j,
	}
}

func (service *ConnectionService) RequestConnection(usernameFrom, usernameTo string) error {
	toUser, err := service.userStore.GetActiveByUsername(usernameTo)
	fromUser, _ := service.userStore.GetActiveByUsername(usernameFrom)
	if err != nil {
		return err
	}
	if toUser.IsPrivate {
		var request = domain.ConnectionRequest{
			Id:          primitive.NewObjectID(),
			From:        *fromUser,
			To:          *toUser,
			RequestTime: time.Now(),
		}
		service.store.Insert(&request)
	} else {
		toUser.Connections = append(toUser.Connections, usernameFrom)
		service.userStore.Update(toUser)
		service.connectionNeo4j.CreateConnectionBetweenUsers(toUser, fromUser)
	}
	fmt.Printf("Saved to db: \n")
	return nil
}

func (service *ConnectionService) GetConnectionUsernamesForUser(username string) ([]string, error) {
	var retVal []string
	connections, _ := service.connectionNeo4j.FindUsersConnection(username)
	for _, connUsername := range connections {
		retVal = append(retVal, connUsername)
	}

	return retVal, nil
}

func (service *ConnectionService) AcceptConnection(connectionId primitive.ObjectID) error {
	connection, err := service.store.Get(connectionId)
	if err != nil {
		return err
	}
	connection.To.Connections = append(connection.To.Connections, connection.From.Username)
	fmt.Printf("Saved connection %s \n", connection.To.Connections)
	service.connectionNeo4j.CreateConnectionBetweenUsers(&connection.From, &connection.To)
	service.store.Delete(connectionId)
	return nil
}

func (service *ConnectionService) DeleteConnection(usernameFrom, usernameTo string) error {
	service.connectionNeo4j.DeleteConnection(usernameFrom, usernameTo)
	return nil
}

func (service *ConnectionService) DeleteConnectionRequest(connectionId primitive.ObjectID) {
	service.store.Delete(connectionId)
}

func (service *ConnectionService) GetRequestsForUser(username string) ([]*domain.ConnectionRequest, error) {
	resp, err := service.store.GetRequestsForUser(username)
	fmt.Printf("Response %d\n", len(resp))
	return resp, err
}

func (service *ConnectionService) GetSuggestedConnectionUsernamesForUser(username string) ([]string, error) {
	var retVal []string
	connections, _ := service.connectionNeo4j.FindSuggestedConnectionsForUser(username)
	for _, connUsername := range connections {
		retVal = append(retVal, connUsername)
	}

	return retVal, nil
}

func (service *ConnectionService) SuggestJobOffersBasedOnUserSkills(username string) ([]*domain.JobOffer, interface{}) {
	var retVal []*domain.JobOffer
	jobOffers, _ := service.connectionNeo4j.FindSuggestedJobOffersBasedOnUserSkills(username)
	for _, jobOffer := range jobOffers {
		retVal = append(retVal, jobOffer)
	}

	return retVal, nil
}
