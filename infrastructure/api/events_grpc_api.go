package api

import (
	"context"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/application"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
)

type EventsHandler struct {
	pb.UnimplementedConnectionServiceServer
	service *application.EventService
}

func NewEventsHandler(service *application.EventService) *EventsHandler {
	return &EventsHandler{
		service: service,
	}
}

func (handler *EventsHandler) GetEvents(ctx context.Context, request *pb.EventRequest) (*pb.Events, error) {
	events, err := handler.service.GetAll()
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
