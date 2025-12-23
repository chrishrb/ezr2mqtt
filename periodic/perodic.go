package periodic

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/store"
	"github.com/chrishrb/ezr2mqtt/transport"
)

type PeriodicRequester struct {
	name     string
	client   transport.Client
	emitter  api.Emitter
	runEvery time.Duration
	store    store.Store
}

func NewPeriodicRequester(
	name string,
	client transport.Client,
	emitter api.Emitter,
	runEvery time.Duration,
	store store.Store,
) *PeriodicRequester {
	return &PeriodicRequester{
		name:     name,
		client:   client,
		emitter:  emitter,
		runEvery: runEvery,
		store:    store,
	}
}

func (r *PeriodicRequester) Run(ctx context.Context) {
	go r.runOnce(ctx)
	go r.run(ctx)
}

func (r *PeriodicRequester) runOnce(ctx context.Context) {
	res, err := r.client.Connect()
	if err != nil {
		slog.Error("error sending periodic message to static endpoint", "error", err)
	}

	// Store device ID
	r.store.SetID(r.name, res.Device.ID)

	// Build json meta data
	metaData := api.Meta{
		Name: res.Device.Name,
	}
	jsonData, err := json.Marshal(metaData)
	if err != nil {
		slog.Error("error marshalling meta data", "error", err)
	}

	// Emit meta data
	id := res.Device.ID
	msg := &api.Message{
		Room: 0, // for complete floor
		Type: "meta",
		Data: jsonData,
	}
	err = r.emitter.Emit(ctx, id, msg)
	if err != nil {
		slog.Error("error emitting meta data message", "error", err)
	}
}

func (r *PeriodicRequester) run(ctx context.Context) {
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

			id := res.Device.ID

			for _, h := range res.Device.HeatAreas {
				msg := &api.Message{
					Room: h.Nr,
					Type: "temperature_target",
					Data: h.TTarget,
				}
				err = r.emitter.Emit(ctx, id, msg)
				if err != nil {
					slog.Error("error emitting periodic temperature_target message", "error", err)
				}

				msg = &api.Message{
					Room: h.Nr,
					Type: "temperature_actual",
					Data: h.TActual,
				}
				err = r.emitter.Emit(ctx, id, msg)
				if err != nil {
					slog.Error("error emitting periodic temperature_actual message", "error", err)
				}
			}
		}
	}
}
