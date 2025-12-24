package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/store"
	"github.com/chrishrb/ezr2mqtt/transport"
)

type HandlerRouter struct {
	client  map[string]transport.Client
	emitter api.Emitter
	store   store.Store
}

func NewHandlerRouter(client map[string]transport.Client, emitter api.Emitter, store store.Store) *HandlerRouter {
	return &HandlerRouter{
		client:  client,
		emitter: emitter,
		store:   store,
	}
}

func (s *HandlerRouter) Handle(ctx context.Context, name string, message *api.Message) {
	var err error

	client, ok := s.client[name]
	if !ok {
		slog.Error("no transport client found for device", "device_name", name)
		return
	}

	id := s.store.GetID(name)
	if id == nil {
		slog.Error("no periodic store ID found for device", "device_name", name)
		return
	}

	err = s.route(client, *id, message)
	if err != nil {
		slog.Error("error handling message", "error", err, "device_name", name, "message_type", message.Type)
	}
}

func (s *HandlerRouter) route(client transport.Client, id string, message *api.Message) error {
	switch message.Type {
	case "temperature_target":
		return setTemperatureTarget(client, id, message)
	case "heatarea_mode":
		return setHeatareaMode(client, id, message)
	default:
		return fmt.Errorf("unknown message type: %s", message.Type)
	}
}
