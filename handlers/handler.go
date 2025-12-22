package handlers

import (
	"context"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/transport"
)

type HandlerRouter struct {
	client transport.Client
}

func NewHandlerRouter(client transport.Client) *HandlerRouter {
	return &HandlerRouter{
		client: client,
	}
}

func (s *HandlerRouter) Handle(ctx context.Context, id string, message *api.Message) {
	// TODO: Implement message handling logic based on message type
}
