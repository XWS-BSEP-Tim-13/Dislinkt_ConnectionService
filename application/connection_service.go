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
		service.connectionNeo4j.CreateConnection(toUser, fromUser)
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
	err1 := service.userStore.Update(&connection.To)
	service.connectionNeo4j.CreateConnection(&connection.From, &connection.To)
	if err != nil {
		return err1
	}
	service.store.Delete(connectionId)
	return nil
}

func (service *ConnectionService) DeleteConnection(usernameFrom, usernameTo string) error {
	user, err := service.userStore.GetActiveByUsername(usernameTo)
	if err != nil {
		return err
	}
	indx := -1
	for i, connection := range user.Connections {
		fmt.Printf("Saved connection %s \n", connection)
		if connection == usernameFrom {
			indx = i
			break
		}
	}
	fmt.Printf("Index %d \n", indx)
	if indx == -1 {
		return nil
	}
	//TODO delete connection between users
	userFrom, err := service.userStore.GetActiveByUsername(usernameFrom)
	service.connectionNeo4j.DeleteConnection(userFrom.Username, user.Username)

	user.Connections[indx] = user.Connections[len(user.Connections)-1]
	user.Connections = user.Connections[:len(user.Connections)-1]
	err = service.userStore.Update(user)
	if err != nil {
		return err
	}
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
