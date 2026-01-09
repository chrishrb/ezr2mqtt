package polling

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
)

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
