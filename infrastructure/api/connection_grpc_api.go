package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/application"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectionHandler struct {
	pb.UnimplementedConnectionServiceServer
	service *application.ConnectionService
}

func NewCompanyHandler(service *application.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{
		service: service,
	}
}

func (handler *ConnectionHandler) GetRequestsForUser(ctx context.Context, request *pb.GetRequestUsername) (*pb.ConnectionRequests, error) {
	//id, err := primitive.ObjectIDFromHex(request.Id)
	//if err != nil {
	//	return nil, err
	//}
	username := request.Username
	requests, _ := handler.service.GetRequestsForUser(username)
	response := &pb.ConnectionRequests{
		Requests: []*pb.ConnectionRequest{},
	}
	for _, request := range requests {
		fmt.Printf("Request: %s, id to: %s\n", request.To.FirstName, request.To.LastName)
		current := mapConnectionRequestToPB(request)
		response.Requests = append(response.Requests, current)
	}
	return response, nil
}

func (handler *ConnectionHandler) AcceptConnectionRequest(ctx context.Context, request *pb.GetRequest) (*pb.ConnectionResponse, error) {
	connectionId, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	handler.service.AcceptConnection(connectionId)
	return new(pb.ConnectionResponse), nil
}

func (handler *ConnectionHandler) DeleteConnectionRequest(ctx context.Context, request *pb.GetRequest) (*pb.ConnectionResponse, error) {
	connectionId, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	handler.service.DeleteConnectionRequest(connectionId)
	return new(pb.ConnectionResponse), nil
}

func (handler *ConnectionHandler) DeleteConnection(ctx context.Context, request *pb.ConnectionBody) (*pb.ConnectionResponse, error) {
	//fmt.Printf("Request: %s, id to: %s\n", request.Connection.IdFrom, request.Connection.IdTo)
	//idFrom, err := primitive.ObjectIDFromHex(request.Connection.IdFrom)
	//idTo, err1 := primitive.ObjectIDFromHex(request.Connection.IdTo)
	//fmt.Printf("Id from: %s, id to: %s\n", idFrom, idTo)
	//if err != nil || err1 != nil {
	//	return nil, err
	//}
	usernameFrom := request.Connection.UsernameFrom
	usernameTo := request.Connection.UsernameTo
	err := handler.service.DeleteConnection(usernameFrom, usernameTo)
	if err != nil {
		return nil, err
	}
	return new(pb.ConnectionResponse), nil
}

func (handler *ConnectionHandler) RequestConnection(ctx context.Context, request *pb.ConnectionBody) (*pb.ConnectionResponse, error) {
	//idFrom, err := primitive.ObjectIDFromHex(request.Connection.IdFrom)
	//idTo, err1 := primitive.ObjectIDFromHex(request.Connection.IdTo)
	//fmt.Printf("Id from: %s, id to: %s\n", idFrom, idTo)
	//if err != nil || err1 != nil {
	//	return nil, err
	//}
	usernameFrom := request.Connection.UsernameFrom
	usernameTo := request.Connection.UsernameTo
	handler.service.RequestConnection(usernameFrom, usernameTo)
	fmt.Printf("Returning to func")
	return new(pb.ConnectionResponse), nil
}

func (handler *ConnectionHandler) GetConnectionUsernamesForUser(ctx context.Context, request *pb.UserUsername) (*pb.UserConnectionUsernames, error) {
	fmt.Printf("Username: %s\n", request.Username)
	connUsernames, err := handler.service.GetConnectionUsernamesForUser(request.Username)
	response := &pb.UserConnectionUsernames{
		Usernames: connUsernames,
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (handler *ConnectionHandler) GetSuggestedConnectionUsernamesForUser(ctx context.Context, request *pb.UserUsername) (*pb.UserConnectionUsernames, error) {
	fmt.Printf("Username: %s\n", request.Username)
	connUsernames, err := handler.service.GetSuggestedConnectionUsernamesForUser(request.Username)
	response := &pb.UserConnectionUsernames{
		Usernames: connUsernames,
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}
