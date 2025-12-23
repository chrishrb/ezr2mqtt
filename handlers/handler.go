package handlers

import (
	"context"
	"fmt"
	"log/slog"

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
	var err error

	switch message.Type {
	case "temperature_target":
		err = s.setTemperatureTarget(id, message)
	}

	if err != nil {
		slog.Error("error handling message", "error", err, "device_id", id, "message_type", message.Type)
	}
}

func (s *HandlerRouter) setTemperatureTarget(id string, message *api.Message) error {
	msg := transport.Message{
		Device: transport.Device{
			ID: id,
			HeatAreas: []transport.HeatArea{{
				Nr:      message.Room,
				TTarget: message.Data.(float64),
			}},
		},
	}
	err := s.client.Send(&msg)
	if err != nil {
		return fmt.Errorf("error sending temperature target: %w", err)
	}
	return nil
}
