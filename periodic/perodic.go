package periodic

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/transport"
)

type PeriodicRequester struct {
	client   transport.Client
	emitter  api.Emitter
	runEvery time.Duration
}

func NewPeriodicRequester(client transport.Client, emitter api.Emitter, runEvery time.Duration) *PeriodicRequester {
	return &PeriodicRequester{
		client:   client,
		emitter:  emitter,
		runEvery: runEvery,
	}
}

func (r *PeriodicRequester) Run(ctx context.Context) {
	go runOnce(ctx, r.client, r.emitter)
	go run(ctx, r.client, r.emitter, r.runEvery)
}

func runOnce(ctx context.Context, client transport.Client, emitter api.Emitter) {
	res, err := client.Connect()
	if err != nil {
		slog.Error("error sending periodic message to static endpoint", "error", err)
	}

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
	err = emitter.Emit(ctx, id, msg)
	if err != nil {
		slog.Error("error emitting meta data message", "error", err)
	}
}

func run(ctx context.Context, client transport.Client, emitter api.Emitter, runEvery time.Duration) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down run periodic")
			return
		case <-time.After(runEvery):
			res, err := client.Connect()
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
				err = emitter.Emit(ctx, id, msg)
				if err != nil {
					slog.Error("error emitting periodic temperature_target message", "error", err)
				}

				msg = &api.Message{
					Room: h.Nr,
					Type: "temperature_actual",
					Data: h.TActual,
				}
				err = emitter.Emit(ctx, id, msg)
				if err != nil {
					slog.Error("error emitting periodic temperature_actual message", "error", err)
				}
			}
		}
	}
}
