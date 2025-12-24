package handlers

import (
	"fmt"
	"strconv"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/transport"
)

func setTemperatureTarget(client transport.Client, id string, message *api.Message) error {
	ttarget, err := strconv.ParseFloat(message.Data, 64)
	if err != nil {
		return fmt.Errorf("invalid temperature target value: %v", message.Data)
	}

	msg := transport.Message{
		Device: transport.Device{
			ID: &id,
			HeatAreas: &[]transport.HeatArea{{
				Nr:      &message.Room,
				TTarget: &ttarget,
			}},
		},
	}

	err = client.Send(&msg)
	if err != nil {
		return fmt.Errorf("error sending temperature target: %w", err)
	}

	return nil
}
