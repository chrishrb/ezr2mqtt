package handlers

import (
	"fmt"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/transport"
)

func setHeatareaMode(client transport.Client, id string, message *api.Message) error {
	var mode int

	switch message.Data.(string) {
	case "auto":
		mode = 0
	case "day":
		mode = 1
	case "night":
		mode = 2
	default:
		return fmt.Errorf("unknown heatarea mode: %s", message.Data.(string))
	}

	msg := transport.Message{
		Device: transport.Device{
			ID: id,
			HeatAreas: []transport.HeatArea{{
				Nr:   message.Room,
				Mode: mode,
			}},
		},
	}

	err := client.Send(&msg)
	if err != nil {
		return fmt.Errorf("error sending heatarea mode: %w", err)
	}

	return nil
}
