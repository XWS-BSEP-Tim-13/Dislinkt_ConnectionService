package application

import (
	events "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/block_user"
	messaging "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/messaging"
)

type BlockUserOrchestrator struct {
	commandPublisher messaging.Publisher
	replySubscriber  messaging.Subscriber
}

func NewBlockUserOrchestrator(publisher messaging.Publisher, subscriber messaging.Subscriber) (*BlockUserOrchestrator, error) {
	o := &BlockUserOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *BlockUserOrchestrator) Start(usernameFrom, usernameTo string) error {
	event := &events.BlockUserCommand{
		Type: events.RemoveConnectionToUser,
		Users: events.Users{
			UserFrom: usernameFrom,
			UserTo:   usernameTo,
		},
	}
	return o.commandPublisher.Publish(event)
}

func (o *BlockUserOrchestrator) handle(reply *events.BlockUserReply) {
	command := events.BlockUserCommand{Users: reply.Users}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *BlockUserOrchestrator) nextCommandType(reply events.BlockUserReplyType) events.BlockUserCommandType {
	switch reply {
	case events.RemoveConnectionFromUserUpdated:
		return events.RemoveConnectionToUser
	case events.RemoveConnectionToUserUpdated:
		return events.BlockUser
	case events.UserBlocked:
		return events.FinnishFunction
	case events.ErrorOccured:
		return events.RollbackUpdates
	default:
		return events.UnknownCommand
	}
}
