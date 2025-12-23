package handlers

import (
	"fmt"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/transport"
)

func setTemperatureTarget(client transport.Client, id string, message *api.Message) error {
	msg := transport.Message{
		Device: transport.Device{
			ID: id,
			HeatAreas: []transport.HeatArea{{
				Nr:      message.Room,
				TTarget: message.Data.(float64),
			}},
		},
	}

	err := client.Send(&msg)
	if err != nil {
		return fmt.Errorf("error sending temperature target: %w", err)
	}
	return nil
}
