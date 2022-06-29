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

func (handler *ConnectionHandler) FindJobOffersBasedOnUserSkills(ctx context.Context, request *pb.UserUsername) (*pb.JobOffers, error) {
	jobOffers, _ := handler.service.SuggestJobOffersBasedOnUserSkills(request.Username)
	response := &pb.JobOffers{
		JobOffers: []*pb.JobOffer{},
	}

	for _, jobOffer := range jobOffers {
		current := mapJobOfferToPB(jobOffer)
		response.JobOffers = append(response.JobOffers, current)
	}

	return response, nil
}
