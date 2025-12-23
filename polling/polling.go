package polling

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/store"
	"github.com/chrishrb/ezr2mqtt/transport"
)

type Poller struct {
	name     string
	client   transport.Client
	emitter  api.Emitter
	runEvery time.Duration
	store    store.Store
}

func NewPoller(
	name string,
	client transport.Client,
	emitter api.Emitter,
	runEvery time.Duration,
	store store.Store,
) *Poller {
	return &Poller{
		name:     name,
		client:   client,
		emitter:  emitter,
		runEvery: runEvery,
		store:    store,
	}
}

func (r *Poller) Run(ctx context.Context) {
	go r.pollOnce(ctx)
	go r.pollPeriodic(ctx)
}

func (r *Poller) pollOnce(ctx context.Context) {
	res, err := r.client.Connect()
	if err != nil {
		slog.Error("error sending message to static endpoint", "error", err)
	}

	// Store device ID
	r.store.SetID(r.name, res.Device.ID)

	// Build json meta data
	rooms := make([]api.RoomDiscovery, len(res.Device.HeatAreas))
	for i, h := range res.Device.HeatAreas {
		rooms[i] = api.RoomDiscovery{
			ID:   h.Nr,
			Name: h.Name,
		}
	}

	metaData := api.ClimateDiscovery{
		// Identity
		Name: res.Device.Name,
		ID:   res.Device.ID,
		Type: res.Device.Type,

		// Rooms
		Rooms: rooms,
	}

	jsonData, err := json.Marshal(metaData)
	if err != nil {
		slog.Error("error marshalling meta data", "error", err)
	}

	// Emit meta data (0 for complete floor)
	r.sendMsg(ctx, 0, "meta", string(jsonData))
}

func (r *Poller) pollPeriodic(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down run periodic")
			return
		case <-time.After(r.runEvery):
			res, err := r.client.Connect()
			if err != nil {
				slog.Error("error sending periodic message to static endpoint", "error", err)
			}

			for _, h := range res.Device.HeatAreas {
				r.sendMsg(ctx, h.Nr, "temperature_target", api.FormatFloat(h.TTarget))
				r.sendMsg(ctx, h.Nr, "temperature_actual", api.FormatFloat(h.TActual))

				mode, err := getHeatAreaMode(h.Mode)
				if err == nil {
					r.sendMsg(ctx, h.Nr, "heatarea_mode", mode)
				} else {
					slog.Error("error getting heat area mode", "error", err)
				}
			}
		}
	}
}

func (r *Poller) sendMsg(ctx context.Context, room int, t string, data string) {
	msg := &api.Message{
		Room: room,
		Type: t,
		Data: data,
	}
	err := r.emitter.Emit(ctx, r.name, msg)
	if err != nil {
		slog.Error("error emitting periodic message", "type", t, "error", err)
	}
}

func getHeatAreaMode(mode int) (string, error) {
	switch mode {
	case 0:
		return "auto", nil
	case 1:
		return "day", nil
	case 2:
		return "night", nil
	default:
		return "", fmt.Errorf("unknown heat area mode: %d", mode)
	}
}
