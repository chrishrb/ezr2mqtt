package handlers

import (
	"testing"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/transport/mock"
	"github.com/stretchr/testify/assert"
)

func TestSetTemperatureTarget_Success(t *testing.T) {
	client := mock.NewMockClient()
	deviceID := "DEVICE-123"

	msg := &api.Message{
		Room: 1,
		Type: "temperature_target",
		Data: 22.5,
	}

	err := setTemperatureTarget(client, deviceID, msg)
	assert.NoError(t, err)

	// Verify the temperature was set in the client
	result, err := client.Connect()
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Check that the device has the correct ID and heat area settings
	assert.Equal(t, deviceID, result.Device.ID)

	// Find the heat area with Nr == 1
	found := false
	for _, heatArea := range result.Device.HeatAreas {
		if heatArea.Nr == 1 {
			assert.Equal(t, 22.5, heatArea.TTarget)
			found = true
			break
		}
	}
	assert.True(t, found, "Heat area 1 should exist and have the correct temperature")
}

func TestSetTemperatureTarget_DifferentRoom(t *testing.T) {
	client := mock.NewMockClient()
	deviceID := "DEVICE-456"

	msg := &api.Message{
		Room: 2,
		Type: "temperature_target",
		Data: 19.5,
	}

	err := setTemperatureTarget(client, deviceID, msg)
	assert.NoError(t, err)

	// Verify the temperature was set in the client
	result, err := client.Connect()
	assert.NoError(t, err)

	// Find the heat area with Nr == 2
	found := false
	for _, heatArea := range result.Device.HeatAreas {
		if heatArea.Nr == 2 {
			assert.Equal(t, 19.5, heatArea.TTarget)
			found = true
			break
		}
	}
	assert.True(t, found, "Heat area 2 should exist and have the correct temperature")
}

func TestSetTemperatureTarget_WithFloatData(t *testing.T) {
	client := mock.NewMockClient()
	deviceID := "DEVICE-789"

	testCases := []struct {
		name        string
		room        int
		temperature float64
	}{
		{"Low temp", 1, 15.0},
		{"Medium temp", 1, 21.0},
		{"High temp", 1, 28.5},
		{"Decimal precision", 2, 22.75},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := &api.Message{
				Room: tc.room,
				Type: "temperature_target",
				Data: tc.temperature,
			}

			err := setTemperatureTarget(client, deviceID, msg)
			assert.NoError(t, err)

			// Verify the temperature was set
			result, err := client.Connect()
			assert.NoError(t, err)

			found := false
			for _, heatArea := range result.Device.HeatAreas {
				if heatArea.Nr == tc.room {
					assert.Equal(t, tc.temperature, heatArea.TTarget)
					found = true
					break
				}
			}
			assert.True(t, found, "Heat area should exist with the correct temperature")
		})
	}
}

func TestSetTemperatureTarget_MessageStructure(t *testing.T) {
	client := mock.NewMockClient()
	deviceID := "TEST-DEVICE"

	msg := &api.Message{
		Room: 3,
		Type: "temperature_target",
		Data: 24.0,
	}

	err := setTemperatureTarget(client, deviceID, msg)
	assert.NoError(t, err)

	// Verify the message was sent by checking the client state
	result, err := client.Connect()
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// The device ID should match
	assert.Equal(t, deviceID, result.Device.ID)

	// Find heat area 3
	found := false
	for _, heatArea := range result.Device.HeatAreas {
		if heatArea.Nr == 3 {
			assert.Equal(t, 24.0, heatArea.TTarget)
			found = true
			break
		}
	}
	assert.True(t, found, "Heat area 3 should have been added with correct temperature")
}
