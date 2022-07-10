package application

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain/enum"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/persistence"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ConnectionService struct {
	store           domain.ConnectionStore
	userStore       domain.UserStore
	connectionNeo4j persistence.ConnectionNeo4jStore
	orchestrator    *BlockUserOrchestrator
	eventStore      domain.EventStore
}

func NewConnectionService(store domain.ConnectionStore, userStore domain.UserStore, neo4j persistence.ConnectionNeo4jStore,
	orchestrator *BlockUserOrchestrator, eventStore domain.EventStore) *ConnectionService {
	return &ConnectionService{
		store:           store,
		userStore:       userStore,
		eventStore:      eventStore,
		connectionNeo4j: neo4j,
		orchestrator:    orchestrator,
	}
}

func (service *ConnectionService) RequestConnection(ctx context.Context, usernameFrom, usernameTo string) (enum.ConnectionStatus, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE RequestConnection")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	toUser, err := service.userStore.GetActiveByUsername(ctx, usernameTo)
	fromUser, _ := service.userStore.GetActiveByUsername(ctx, usernameFrom)

	ret := enum.CONNECTED
	fmt.Println("Request", usernameTo, usernameFrom)
	if err != nil {
		return enum.NONE, err
	}
	if toUser.IsPrivate {
		fmt.Println("Is private")
		var request = domain.ConnectionRequest{
			Id:          primitive.NewObjectID(),
			From:        *fromUser,
			To:          *toUser,
			RequestTime: time.Now(),
		}
		ret = enum.CONNECTION_REQUEST
		service.store.Insert(ctx, &request)
	} else {
		fmt.Println("Is not")
		fromUser.Connections = append(fromUser.Connections, usernameTo)
		service.userStore.Update(ctx, fromUser)
		fmt.Println(fromUser)
		service.connectionNeo4j.CreateConnectionBetweenUsers(toUser, fromUser)
	}
	var event = domain.Event{Id: primitive.NewObjectID(), User: usernameFrom, Action: `Create connection request to user ` + usernameTo, Published: time.Now()}
	service.eventStore.Insert(&event)

	fmt.Printf("Saved to db: \n")
	return ret, nil
}

func (service *ConnectionService) BlockOrchestrator(ctx context.Context, usernameFrom, usernameTo string) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE BlockOrchestrator")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	fmt.Println("Orchestrator stared...")
	err := service.orchestrator.Start(usernameFrom, usernameTo)
	if err != nil {
		return err
	}
	return nil
}

func (service *ConnectionService) GetConnectionUsernamesForUser(ctx context.Context, username string) ([]string, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE GetConnectionUsernamesForUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	var retVal []string
	connections, _ := service.connectionNeo4j.FindUsersConnection(username)
	for _, connUsername := range connections {
		retVal = append(retVal, connUsername)
	}

	return retVal, nil
}

func (service *ConnectionService) GetConnectionSuggestionsForUser(ctx context.Context, username string) ([]string, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE GetConnectionSuggestionsForUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	var retVal []string
	connections, _ := service.connectionNeo4j.FindSuggestedConnectionsForUser(username)
	for _, connUsername := range connections {
		retVal = append(retVal, connUsername)
	}

	return retVal, nil
}

func (service *ConnectionService) AcceptConnection(ctx context.Context, usernameFrom, usernameTo string) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE AcceptConnection")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	connection, err := service.store.GetConnectionByUsernames(ctx, usernameFrom, usernameTo)
	if err != nil {
		return err
	}
	user, err := service.userStore.GetByUsername(ctx, usernameTo)
	if err != nil {
		return err
	}
	user.Connections = append(user.Connections, usernameFrom)
	connection.To.Connections = append(connection.To.Connections, connection.From.Username)
	fmt.Printf("Saved connection %s \n", connection.To.Connections)
	service.connectionNeo4j.CreateConnectionBetweenUsers(&connection.From, &connection.To)
	service.connectionNeo4j.CreateConnectionBetweenUsers(&connection.To, &connection.From)
	service.store.Delete(ctx, connection.Id)
	var event = domain.Event{Id: primitive.NewObjectID(), User: usernameTo, Action: `Accepted connection request from user ` + usernameFrom, Published: time.Now()}
	service.eventStore.Insert(&event)
	return nil
}

func (service *ConnectionService) DeleteConnection(ctx context.Context, usernameFrom, usernameTo string) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE DeleteConnection")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := service.userStore.GetActiveByUsername(ctx, usernameTo)
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
	err = service.userStore.Update(ctx, user)
	fmt.Println("Delete stared...", usernameTo)
	service.connectionNeo4j.DeleteConnection(usernameFrom, usernameTo)
	var event = domain.Event{Id: primitive.NewObjectID(), User: usernameTo, Action: `Delete connection ` + usernameFrom, Published: time.Now()}
	service.eventStore.Insert(&event)
	return nil
}

func (service *ConnectionService) UnblockUser(ctx context.Context, usernameFrom, usernameTo string) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE UnblockUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := service.userStore.GetActiveByUsername(ctx, usernameTo)
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

	err = service.userStore.Update(ctx, user)
	var event = domain.Event{Id: primitive.NewObjectID(), User: usernameFrom, Action: `Unblocked user ` + usernameTo, Published: time.Now()}
	service.eventStore.Insert(&event)

	return nil
}

func (service *ConnectionService) DeleteConnectionRequest(ctx context.Context, usernameFrom, usernameTo string) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE DeleteConnectionRequest")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	fmt.Println("delete connection request", usernameTo, usernameFrom)

	request, _ := service.store.GetConnectionByUsernames(ctx, usernameFrom, usernameTo)
	service.store.Delete(ctx, request.Id)
	var event = domain.Event{Id: primitive.NewObjectID(), User: usernameTo, Action: `Delete connection request from  ` + usernameFrom, Published: time.Now()}
	service.eventStore.Insert(&event)

}

func (service *ConnectionService) GetRequestsForUser(ctx context.Context, username string) ([]*domain.ConnectionRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE GetRequestsForUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	resp, err := service.store.GetRequestsForUser(ctx, username)
	fmt.Printf("Response %d\n", len(resp))
	return resp, err
}

func (service *ConnectionService) BlockUser(ctx context.Context, usernameFrom, usernameTo string) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE BlockUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	fmt.Println("Block user")
	user, err := service.userStore.GetByUsername(ctx, usernameFrom)
	if err != nil {
		return err
	}
	user.BlockedUsers = append(user.BlockedUsers, usernameTo)

	var event = domain.Event{Id: primitive.NewObjectID(), User: usernameFrom, Action: `Blocked user ` + usernameTo, Published: time.Now()}
	service.eventStore.Insert(&event)
	err = service.userStore.UpdateBlockedList(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (service *ConnectionService) CheckIfUserConnected(ctx context.Context, fromUsername, toUsername string) enum.ConnectionStatus {
	span := tracer.StartSpanFromContext(ctx, "SERVICE CheckIfUserConnected")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	resp, err := service.userStore.CheckIfUserIsBlocked(ctx, toUsername, fromUsername)
	fmt.Println(resp, err)
	if resp != nil {
		return enum.BLOCKED
	}
	resp, err = service.userStore.CheckIfUserIsBlocked(ctx, fromUsername, toUsername)
	fmt.Println(resp, err)
	if resp != nil {
		return enum.BLOCKED_ME
	}
	resp, err = service.userStore.CheckIfUsersConnected(ctx, toUsername, fromUsername)
	fmt.Println(resp, err)
	if err == nil {
		return enum.CONNECTED
	}
	isRequested := service.store.CheckIfUsersConnected(ctx, fromUsername, toUsername)
	fmt.Println(isRequested)
	if isRequested {
		return enum.CONNECTION_REQUEST
	}
	return enum.NONE
}

func (service *ConnectionService) SuggestJobOffersBasedOnUserSkills(ctx context.Context, username string) ([]*domain.JobOffer, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE SuggestJobOffersBasedOnUserSkills")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	var retVal []*domain.JobOffer
	jobOffers, _ := service.connectionNeo4j.FindSuggestedJobOffersBasedOnUserSkills(username)
	for _, jobOffer := range jobOffers {
		retVal = append(retVal, jobOffer)
	}

	return retVal, nil
}

func (service *ConnectionService) GetAllEvents() ([]*domain.Event, error) {
	return service.eventStore.GetAll()
}
