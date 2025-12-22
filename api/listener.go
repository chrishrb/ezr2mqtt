package api

import "context"

type MessageHandler interface {
	Handle(ctx context.Context, id string, message *Message)
}

type MessageHandlerFunc func(ctx context.Context, id string, message *Message)

func (h MessageHandlerFunc) Handle(ctx context.Context, id string, message *Message) {
	h(ctx, id, message)
}

type Listener interface {
	Connect(ctx context.Context, handler MessageHandler) (Connection, error)
}

type Connection interface {
	Disconnect(ctx context.Context) error
}
