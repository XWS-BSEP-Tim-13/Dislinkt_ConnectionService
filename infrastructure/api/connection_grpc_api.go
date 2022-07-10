package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/application"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/jwt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "API GetRequestsForUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	username := request.Username
	requests, _ := handler.service.GetRequestsForUser(ctx, username)
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
	span := tracer.StartSpanFromContext(ctx, "API AcceptConnectionRequest")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Delete con request started")
	handler.service.AcceptConnection(ctx, request.Username, username)
	return new(pb.ConnectionResponse), nil
}

func (handler *ConnectionHandler) DeleteConnectionRequest(ctx context.Context, request *pb.UserUsername) (*pb.ConnectionResponse, error) {
	username, err := jwt.ExtractUsernameFromToken(ctx)
	span := tracer.StartSpanFromContext(ctx, "API DeleteConnectionRequest")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Delete con request started")
	handler.service.DeleteConnectionRequest(ctx, request.Username, username)
	return new(pb.ConnectionResponse), nil
}

func (handler *ConnectionHandler) DeleteConnection(ctx context.Context, request *pb.UserUsername) (*pb.ConnectionResponse, error) {
	username, err := jwt.ExtractUsernameFromToken(ctx)
	span := tracer.StartSpanFromContext(ctx, "API DeleteConnection")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = handler.service.DeleteConnection(ctx, request.Username, username)
	if err != nil {
		return nil, err
	}
	return new(pb.ConnectionResponse), nil
}

func (handler *ConnectionHandler) RequestConnection(ctx context.Context, request *pb.UserUsername) (*pb.ConnectionStatusResponse, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	span := tracer.StartSpanFromContext(ctx, "API RequestConnection")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	fmt.Println(request.Username)
	ret, err := handler.service.RequestConnection(ctx, username, request.Username)
	if err != nil {
		return nil, err
	}
	pbVal := pb.ConnectionStatus(ret)

	response := &pb.ConnectionStatusResponse{
		ConnectionStatus: pbVal,
	}
	return response, nil
}

func (handler *ConnectionHandler) GetConnectionUsernamesForUser(ctx context.Context, request *pb.ConnectionResponse) (*pb.UserConnectionUsernames, error) {
	username, err := jwt.ExtractUsernameFromToken(ctx)
	span := tracer.StartSpanFromContext(ctx, "API GetConnectionUsernamesForUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	if err != nil {
		response := &pb.UserConnectionUsernames{
			Usernames: []string{},
		}
		return response, nil
	}
	fmt.Println(username)
	connUsernames, err := handler.service.GetConnectionUsernamesForUser(ctx, username)
	response := &pb.UserConnectionUsernames{
		Usernames: connUsernames,
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (handler *ConnectionHandler) GetSuggestedConnectionUsernamesForUser(ctx context.Context, request *pb.UserUsername) (*pb.UserConnectionUsernames, error) {
	username, err := jwt.ExtractUsernameFromToken(ctx)
	span := tracer.StartSpanFromContext(ctx, "API GetSuggestedConnectionUsernamesForUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	if err != nil {
		response := &pb.UserConnectionUsernames{
			Usernames: []string{},
		}
		return response, nil
	}
	fmt.Println(username)
	connUsernames, err := handler.service.GetConnectionSuggestionsForUser(ctx, username)
	response := &pb.UserConnectionUsernames{
		Usernames: connUsernames,
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (handler *ConnectionHandler) FindJobOffersBasedOnUserSkills(ctx context.Context, request *pb.UserUsername) (*pb.JobOffers, error) {
	username, _ := jwt.ExtractUsernameFromToken(ctx)
	span := tracer.StartSpanFromContext(ctx, "API FindJobOffersBasedOnUserSkills")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	jobs, err := handler.service.SuggestJobOffersBasedOnUserSkills(ctx, username)

	if err != nil {
		return nil, err
	}

	response := &pb.JobOffers{
		JobOffers: []*pb.JobOffer{},
	}

	for _, offer := range jobs {
		current := mapJobOfferToPB(offer)
		response.JobOffers = append(response.JobOffers, current)
	}

	return response, nil
}

func (handler *ConnectionHandler) CreateJobOffer(ctx context.Context, request *pb.JobOfferRequest) (*pb.GetRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "API CreateJobOffer")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	job := mapJobOfferDtoToDomain(request.Dto)
	err := handler.service.InsertJobOffer(ctx, job)
	if err != nil {
		return nil, err
	}

	response := &pb.GetRequest{}
	return response, nil
}

func (handler *ConnectionHandler) UnBlockUser(ctx context.Context, request *pb.UserUsername) (*pb.GetAllRequest, error) {
	usernameTo, err := jwt.ExtractUsernameFromToken(ctx)
	span := tracer.StartSpanFromContext(ctx, "API UnBlockUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	err = handler.service.UnblockUser(ctx, request.Username, usernameTo)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := &pb.GetAllRequest{}
	return response, nil
}

func (handler *ConnectionHandler) BlockUser(ctx context.Context, request *pb.UserUsername) (*pb.GetAllRequest, error) {
	usernameFrom, err := jwt.ExtractUsernameFromToken(ctx)
	span := tracer.StartSpanFromContext(ctx, "API BlockUser")
	defer span.Finish()
	fmt.Println(usernameFrom)
	ctx = tracer.ContextWithSpan(context.Background(), span)

	fmt.Println("block started")
	if err != nil {
		fmt.Println("Stopped on extracting username")
		fmt.Println(err)
		return nil, err
	}
	err = handler.service.BlockOrchestrator(ctx, usernameFrom, request.Username)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := &pb.GetAllRequest{}
	return response, nil
}

func (handler *ConnectionHandler) CheckIfUserConnected(ctx context.Context, request *pb.UserUsername) (*pb.ConnectionStatusResponse, error) {
	usernameFrom, err := jwt.ExtractUsernameFromToken(ctx)
	span := tracer.StartSpanFromContext(ctx, "API CheckIfUserConnected")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	fmt.Println("Check has started", usernameFrom)
	if err != nil {
		fmt.Println(err)
		response := &pb.ConnectionStatusResponse{
			ConnectionStatus: 4,
		}
		return response, nil
	}
	isConnected := handler.service.CheckIfUserConnected(ctx, usernameFrom, request.Username)
	fmt.Println(isConnected)
	pbVal := pb.ConnectionStatus(isConnected)

	response := &pb.ConnectionStatusResponse{
		ConnectionStatus: pbVal,
	}
	return response, nil
}

func (handler *ConnectionHandler) GetEvents(ctx context.Context, request *pb.EventRequest) (*pb.Events, error) {
	events, err := handler.service.GetAllEvents()
	if err != nil {
		return nil, err
	}
	response := &pb.Events{
		Events: []*pb.Event{},
	}

	for _, event := range events {
		current := mapEventToPB(event)
		response.Events = append(response.Events, current)
	}

	return response, nil
}
