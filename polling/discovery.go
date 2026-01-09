package polling

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/transport"
)

func (r *Poller) runDiscovery(ctx context.Context) {
	// Run discovery immediately at startup
	r.doDiscovery(ctx)

	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down home assistant discovery")
			return
		case <-time.After(r.discoveryRunEvery):
			r.doDiscovery(ctx)
		}
	}
}

func (r *Poller) doDiscovery(ctx context.Context) {
	res, err := r.client.Connect()
	if err != nil {
		slog.Error("error sending message to static endpoint", "error", err)
		return
	}

	// Store device ID
	r.store.SetID(r.name, *res.Device.ID)

	// Build json meta data
	if res.Device.HeatAreas != nil {
		for _, h := range *res.Device.HeatAreas {
			roomName := removeUmlauts(*h.Name)
			roomNumber := *h.Nr

			r.emitTemperatureTargetDiscovery(ctx, res, roomName, roomNumber, *h.TTargetMin, *h.TTargetMax)
			r.emitTemperatureActualDiscovery(ctx, res, roomName, roomNumber)
			r.emitHeatareaModeDiscovery(ctx, res, roomName, roomNumber)
		}
	}
}

func (r *Poller) emitTemperatureTargetDiscovery(ctx context.Context, res *transport.Message, roomName string, roomNumber int, targetMin, targetMax float64) {
	id := fmt.Sprintf("%s-%s-temperature_target", r.name, strings.ToLower(roomName))
	err := r.emitter.EmitHADiscovery(ctx, api.HAComponentNumber, &api.HADiscovery{
		Name:              id,
		UniqueID:          id,
		StateTopic:        fmt.Sprintf("%s/%s/%d/state/temperature_target", r.emitter.GetPrefix(), r.name, roomNumber),
		UnitOfMeasurement: "°C",
		DeviceClass:       "temperature",
		StateClass:        "measurement",
		CommandTopic:      fmt.Sprintf("%s/%s/%d/set/temperature_target", r.emitter.GetPrefix(), r.name, roomNumber),
		Minimum:           targetMin,
		Maximum:           targetMax,
		Step:              0.5,
		Mode:              "slider",
		Device: &api.HADevice{
			Identifiers: []string{*res.Device.ID},
			Name:        *res.Device.Name,
		},
	})

	if err != nil {
		slog.Error("error emitting HA discovery for temperature target", "error", err)
	}
}

func (r *Poller) emitTemperatureActualDiscovery(ctx context.Context, res *transport.Message, roomName string, roomNumber int) {
	id := fmt.Sprintf("%s-%s-temperature_actual", r.name, strings.ToLower(roomName))
	err := r.emitter.EmitHADiscovery(ctx, api.HAComponentSensor, &api.HADiscovery{
		Name:              id,
		UniqueID:          id,
		StateTopic:        fmt.Sprintf("%s/%s/%d/state/temperature_actual", r.emitter.GetPrefix(), r.name, roomNumber),
		UnitOfMeasurement: "°C",
		DeviceClass:       "temperature",
		StateClass:        "measurement",
		Device: &api.HADevice{
			Identifiers: []string{*res.Device.ID},
			Name:        *res.Device.Name,
		},
	})
	if err != nil {
		slog.Error("error emitting HA discovery for temperature actual", "error", err)
	}
}

func (r *Poller) emitHeatareaModeDiscovery(ctx context.Context, res *transport.Message, roomName string, roomNumber int) {
	id := fmt.Sprintf("%s-%s-heatarea_mode", r.name, strings.ToLower(roomName))
	err := r.emitter.EmitHADiscovery(ctx, api.HAComponentSelect, &api.HADiscovery{
		Name:         id,
		UniqueID:     id,
		StateTopic:   fmt.Sprintf("%s/%s/%d/state/heatarea_mode", r.emitter.GetPrefix(), r.name, roomNumber),
		CommandTopic: fmt.Sprintf("%s/%s/%d/set/heatarea_mode", r.emitter.GetPrefix(), r.name, roomNumber),
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
	if err != nil {
		slog.Error("error emitting HA discovery for heatarea mode", "error", err)
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
