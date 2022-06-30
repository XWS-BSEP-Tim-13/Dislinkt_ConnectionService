package application

import (
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain/enum"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ConnectionService struct {
	store           domain.ConnectionStore
	userStore       domain.UserStore
	connectionNeo4j persistence.ConnectionNeo4jStore
	orchestrator    *BlockUserOrchestrator
}

func NewConnectionService(store domain.ConnectionStore, userStore domain.UserStore, neo4j persistence.ConnectionNeo4jStore, orchestrator *BlockUserOrchestrator) *ConnectionService {
	return &ConnectionService{
		store:           store,
		userStore:       userStore,
		connectionNeo4j: neo4j,
		orchestrator:    orchestrator,
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

func (service *ConnectionService) BlockOrchestrator(usernameFrom, usernameTo string) error {
	fmt.Println("Orchestrator stared...")
	err := service.orchestrator.Start(usernameFrom, usernameTo)
	if err != nil {
		return err
	}
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

func (service *ConnectionService) AcceptConnection(usernameFrom, usernameTo string) error {
	connection, err := service.store.GetConnectionByUsernames(usernameFrom, usernameTo)
	if err != nil {
		return err
	}
	user, err := service.userStore.GetByUsername(usernameTo)
	if err != nil {
		return err
	}
	user.Connections = append(user.Connections, usernameFrom)
	connection.To.Connections = append(connection.To.Connections, connection.From.Username)
	fmt.Printf("Saved connection %s \n", connection.To.Connections)
	service.connectionNeo4j.CreateConnectionBetweenUsers(&connection.From, &connection.To)
	service.store.Delete(connection.Id)
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
	user.Connections[indx] = user.Connections[len(user.Connections)-1]
	user.Connections = user.Connections[:len(user.Connections)-1]
	err = service.userStore.Update(user)
	fmt.Println("Delete stared...", usernameTo)
	service.connectionNeo4j.DeleteConnection(usernameFrom, usernameTo)
	return nil
}

func (service *ConnectionService) UnblockUser(usernameFrom, usernameTo string) error {
	user, err := service.userStore.GetActiveByUsername(usernameTo)
	if err != nil {
		return err
	}
	indx := -1
	for i, blocks := range user.BlockedUsers {
		fmt.Printf("Saved connection %s \n", blocks)
		if blocks == usernameFrom {
			indx = i
			break
		}
	}
	fmt.Printf("Index %d \n", indx)
	if indx == -1 {
		return nil
	}
	user.BlockedUsers[indx] = user.BlockedUsers[len(user.BlockedUsers)-1]
	user.BlockedUsers = user.BlockedUsers[:len(user.BlockedUsers)-1]
	err = service.userStore.Update(user)
	return nil
}

func (service *ConnectionService) DeleteConnectionRequest(usernameFrom, usernameTo string) {
	fmt.Println("delete connection request", usernameTo, usernameFrom)
	request, _ := service.store.GetConnectionByUsernames(usernameFrom, usernameTo)
	service.store.Delete(request.Id)
}

func (service *ConnectionService) GetRequestsForUser(username string) ([]*domain.ConnectionRequest, error) {
	resp, err := service.store.GetRequestsForUser(username)
	fmt.Printf("Response %d\n", len(resp))
	return resp, err
}

func (service *ConnectionService) BlockUser(usernameFrom, usernameTo string) error {
	fmt.Println("Block user")
	user, err := service.userStore.GetByUsername(usernameFrom)
	if err != nil {
		return err
	}
	user.BlockedUsers = append(user.BlockedUsers, usernameTo)
	err = service.userStore.UpdateBlockedList(user)
	if err != nil {
		return err
	}
	return nil
}

func (service *ConnectionService) CheckIfUserConnected(fromUsername, toUsername string) enum.ConnectionStatus {
	resp, err := service.userStore.CheckIfUserIsBlocked(toUsername, fromUsername)
	fmt.Println(resp, err)
	if resp != nil {
		return enum.BLOCKED
	}
	resp, err = service.userStore.CheckIfUserIsBlocked(fromUsername, toUsername)
	fmt.Println(resp, err)
	if resp != nil {
		return enum.BLOCKED_ME
	}
	resp, err = service.userStore.CheckIfUsersConnected(fromUsername, toUsername)
	fmt.Println(resp, err)
	if err == nil {
		return enum.CONNECTED
	}
	isRequested := service.store.CheckIfUsersConnected(fromUsername, toUsername)
	fmt.Println(isRequested)
	if isRequested {
		return enum.CONNECTION_REQUEST
	}
	return enum.NONE
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
