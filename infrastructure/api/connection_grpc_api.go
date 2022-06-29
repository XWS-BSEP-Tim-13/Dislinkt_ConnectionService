package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/application"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/jwt"
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
	usernameFrom := request.Connection.UsernameFrom
	usernameTo := request.Connection.UsernameTo
	err := handler.service.DeleteConnection(usernameFrom, usernameTo)
	if err != nil {
		return nil, err
	}
	return new(pb.ConnectionResponse), nil
}

func (handler *ConnectionHandler) RequestConnection(ctx context.Context, request *pb.ConnectionBody) (*pb.ConnectionResponse, error) {
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

func (handler *ConnectionHandler) BlockUser(ctx context.Context, request *pb.UserUsername) (*pb.GetAllRequest, error) {
	usernameFrom, err := jwt.ExtractUsernameFromToken(ctx)
	if err != nil {
		return nil, err
	}
	err = handler.service.BlockOrchestrator(usernameFrom, request.Username)
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllRequest{}
	return response, nil
}

func (handler *ConnectionHandler) CheckIfUserConnected(ctx context.Context, request *pb.UserUsername) (*pb.ConnectionStatusResponse, error) {
	usernameFrom, err := jwt.ExtractUsernameFromToken(ctx)
	if err != nil {
		fmt.Println(err)
		response := &pb.ConnectionStatusResponse{
			ConnectionStatus: 4,
		}
		return response, nil
	}
	isConnected := handler.service.CheckIfUserConnected(usernameFrom, request.Username)
	fmt.Println(isConnected)
	pbVal := pb.ConnectionStatus(isConnected)

	response := &pb.ConnectionStatusResponse{
		ConnectionStatus: pbVal,
	}
	return response, nil
}
