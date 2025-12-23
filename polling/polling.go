package polling

import (
	"context"
	"encoding/json"
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

	// Emit meta data
	msg := &api.Message{
		Room: 0, // for complete floor
		Type: "meta",
		Data: jsonData,
	}
	err = r.emitter.Emit(ctx, r.name, msg)
	if err != nil {
		slog.Error("error emitting meta data message", "error", err)
	}
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
				msg := &api.Message{
					Room: h.Nr,
					Type: "temperature_target",
					Data: h.TTarget,
				}
				err = r.emitter.Emit(ctx, r.name, msg)
				if err != nil {
					slog.Error("error emitting periodic temperature_target message", "error", err)
				}

				msg = &api.Message{
					Room: h.Nr,
					Type: "temperature_actual",
					Data: h.TActual,
				}
				err = r.emitter.Emit(ctx, r.name, msg)
				if err != nil {
					slog.Error("error emitting periodic temperature_actual message", "error", err)
				}
			}
		}
	}
}
