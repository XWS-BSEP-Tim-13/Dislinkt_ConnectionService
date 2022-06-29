package api

import (
	events "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/block_user"
	saga "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/messaging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/application"
)

type BlockUserCommandHandler struct {
	connectionService *application.ConnectionService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewBlockUserCommandHandler(orderService *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) (*BlockUserCommandHandler, error) {
	o := &BlockUserCommandHandler{
		connectionService: orderService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *BlockUserCommandHandler) handle(command *events.BlockUserCommand) {
	reply := events.BlockUserReply{Users: command.Users}

	switch command.Type {
	case events.RemoveConnectionFromUser:
		err := handler.connectionService.DeleteConnection(reply.Users.UserFrom, reply.Users.UserTo)
		if err != nil {
			return
		}
		reply.Type = events.RemoveConnectionFromUserUpdated
	case events.RemoveConnectionToUser:
		err := handler.connectionService.DeleteConnection(reply.Users.UserTo, reply.Users.UserFrom)
		if err != nil {
			return
		}
		reply.Type = events.RemoveConnectionFromUserUpdated
	case events.BlockUser:
		err := handler.connectionService.BlockUser(reply.Users.UserFrom, reply.Users.UserTo)
		if err != nil {
			return
		}
	case events.FinnishFunction:
		return
	case events.RollbackUpdates:
		return
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
