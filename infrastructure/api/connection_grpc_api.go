package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/application"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/jwt"
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

func (handler *ConnectionHandler) AcceptConnectionRequest(ctx context.Context, request *pb.UserUsername) (*pb.ConnectionResponse, error) {
	username, err := jwt.ExtractUsernameFromToken(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Delete con request started")
	handler.service.AcceptConnection(request.Username, username)
	return new(pb.ConnectionResponse), nil
}

func (handler *ConnectionHandler) DeleteConnectionRequest(ctx context.Context, request *pb.UserUsername) (*pb.ConnectionResponse, error) {
	username, err := jwt.ExtractUsernameFromToken(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Delete con request started")
	handler.service.DeleteConnectionRequest(username, request.Username)
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

func (handler *ConnectionHandler) RequestConnection(ctx context.Context, request *pb.UserUsername) (*pb.ConnectionStatusResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	fmt.Println(request.Username)
	ret, err := handler.service.RequestConnection(username, request.Username)
	if err != nil {
		return nil, err
	}
	pbVal := pb.ConnectionStatus(ret)

	response := &pb.ConnectionStatusResponse{
		ConnectionStatus: pbVal,
	}
	return response, nil
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

func (handler *ConnectionHandler) UnBlockUser(ctx context.Context, request *pb.UserUsername) (*pb.GetAllRequest, error) {
	usernameTo, err := jwt.ExtractUsernameFromToken(ctx)
	err = handler.service.UnblockUser(request.Username, usernameTo)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := &pb.GetAllRequest{}
	return response, nil
}

func (handler *ConnectionHandler) BlockUser(ctx context.Context, request *pb.UserUsername) (*pb.GetAllRequest, error) {
	fmt.Println("block started")
	usernameFrom, err := jwt.ExtractUsernameFromToken(ctx)
	if err != nil {
		fmt.Println("Stopped on extracting username")
		fmt.Println(err)
		return nil, err
	}
	err = handler.service.BlockOrchestrator(usernameFrom, request.Username)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := &pb.GetAllRequest{}
	return response, nil
}

func (handler *ConnectionHandler) CheckIfUserConnected(ctx context.Context, request *pb.UserUsername) (*pb.ConnectionStatusResponse, error) {
	usernameFrom, err := jwt.ExtractUsernameFromToken(ctx)
	fmt.Println("Check has started", usernameFrom)
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
