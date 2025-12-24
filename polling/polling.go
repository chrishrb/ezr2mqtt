package polling

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
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
	r.store.SetID(r.name, *res.Device.ID)

	// Build json meta data
	if res.Device.HeatAreas != nil {
		for _, h := range *res.Device.HeatAreas {
			roomName := removeUmlauts(*h.Name)
			roomNumber := *h.Nr

			r.emitter.EmitHADiscovery(ctx, api.HAComponentNumber, api.HASensorDiscovery{
				Name:     fmt.Sprintf("%s Temperature Target", roomName),
				UniqueID: fmt.Sprintf("%s-%s-temperature_target", r.name, strings.ToLower(roomName)),
				// TODO: refactor
				StateTopic:        fmt.Sprintf("%s/%s/%d/state/temperature_target", "ezr", r.name, roomNumber),
				UnitOfMeasurement: "°C",
				DeviceClass:       "temperature",
				StateClass:        "measurement",
				// TODO: refactor
				CommandTopic: fmt.Sprintf("%s/%s/%d/set/temperature_target", "ezr", r.name, roomNumber),
				Minimum:      *h.TTargetMin,
				Maximum:      *h.TTargetMax,
				Step:         0.5,
				Mode:         "slider",
				Device: &api.HADevice{
					Identifiers: []string{*res.Device.ID},
					Name:        *res.Device.Name,
				},
			})

			r.emitter.EmitHADiscovery(ctx, api.HAComponentSensor, api.HASensorDiscovery{
				Name:     fmt.Sprintf("%s Temperature Actual", roomName),
				UniqueID: fmt.Sprintf("%s-%s-temperature_actual", r.name, strings.ToLower(roomName)),
				// TODO: refactor
				StateTopic:        fmt.Sprintf("%s/%s/%d/state/temperature_actual", "ezr", r.name, roomNumber),
				UnitOfMeasurement: "°C",
				DeviceClass:       "temperature",
				StateClass:        "measurement",
				Device: &api.HADevice{
					Identifiers: []string{*res.Device.ID},
					Name:        *res.Device.Name,
				},
			})

			r.emitter.EmitHADiscovery(ctx, api.HAComponentSelect, api.HASensorDiscovery{
				Name:     fmt.Sprintf("%s Heatarea Mode", roomName),
				UniqueID: fmt.Sprintf("%s-%s-heatarea_mode", r.name, strings.ToLower(roomName)),
				// TODO: refactor
				StateTopic: fmt.Sprintf("%s/%s/%d/state/heatarea_mode", "ezr", r.name, roomNumber),
				// TODO: refactor
				CommandTopic: fmt.Sprintf("%s/%s/%d/set/heatarea_mode", "ezr", r.name, roomNumber),
				Options: []string{
					"auto",
					"day",
					"night",
				},
				Device: &api.HADevice{
					Identifiers: []string{*res.Device.ID},
					Name:        *res.Device.Name,
				},
			})
		}
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
				continue
			}

			if res.Device.HeatAreas != nil {
				for _, h := range *res.Device.HeatAreas {
					roomNumber := *h.Nr

					r.sendMsg(ctx, roomNumber, "temperature_target", api.FormatFloat(*h.TTarget))
					r.sendMsg(ctx, roomNumber, "temperature_actual", api.FormatFloat(*h.TActual))

					mode, err := getHeatAreaMode(*h.Mode)
					if err == nil {
						r.sendMsg(ctx, roomNumber, "heatarea_mode", mode)
					} else {
						slog.Error("error getting heat area mode", "error", err)
					}
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
func removeUmlauts(s string) string {
	s = strings.ReplaceAll(s, "ä", "ae")
	s = strings.ReplaceAll(s, "ö", "oe")
	s = strings.ReplaceAll(s, "ü", "ue")
	s = strings.ReplaceAll(s, "ß", "ss")
	s = strings.ReplaceAll(s, "Ä", "Ae")
	s = strings.ReplaceAll(s, "Ö", "Oe")
	s = strings.ReplaceAll(s, "Ü", "Ue")
	return s
}
