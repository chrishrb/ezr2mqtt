package handlers

import (
	"context"
	"testing"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/store"
	"github.com/chrishrb/ezr2mqtt/transport"
	"github.com/chrishrb/ezr2mqtt/transport/mock"
	"github.com/stretchr/testify/assert"
)

type mockEmitter struct {
	EmitWasCalled            bool
	EmitHADiscoveryWasCalled bool
}

func (m *mockEmitter) Emit(ctx context.Context, id string, message *api.Message) error {
	m.EmitWasCalled = true
	return nil
}

func (m *mockEmitter) EmitHADiscovery(ctx context.Context, component api.HAComponent, message *api.HADiscovery) error {
	m.EmitHADiscoveryWasCalled = true
	return nil
}

func (m *mockEmitter) GetPrefix() string {
	return "ezr"
}

func TestNewHandlerRouter(t *testing.T) {
	client := mock.NewMockClient()
	emitter := &mockEmitter{}
	store := store.NewInMemoryStore()
	clientMap := map[string]transport.Client{
		"device1": client,
	}

	router := NewHandlerRouter(clientMap, emitter, store)

	assert.NotNil(t, router)
	assert.Equal(t, clientMap, router.client)
	assert.Equal(t, store, router.store)
}

func TestHandlerRouter_Handle_Success(t *testing.T) {
	client := mock.NewMockClient()
	emitter := &mockEmitter{}
	store := store.NewInMemoryStore()
	deviceName := "device1"
	deviceID := "DEVICE-123"

	// Setup store with device ID
	store.SetID(deviceName, deviceID)

	clientMap := map[string]transport.Client{
		deviceName: client,
	}

	router := NewHandlerRouter(clientMap, emitter, store)

	msg := &api.Message{
		Room: 1,
		Type: "temperature_target",
		Data: "21.5",
	}

	ctx := context.Background()
	router.Handle(ctx, deviceName, msg)

	// Verify that the message was sent to the client
	result, err := client.Connect()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Device.HeatAreas)

	// Verify the temperature was set correctly
	found := false
	for _, heatArea := range *result.Device.HeatAreas {
		if *heatArea.Nr == 1 {
			assert.Equal(t, 21.5, *heatArea.TTarget)
			found = true
			break
		}
	}
	assert.True(t, found, "Heat area 1 should exist")
}

func TestHandlerRouter_Handle_NoClient(t *testing.T) {
	store := store.NewInMemoryStore()
	emitter := &mockEmitter{}
	clientMap := map[string]transport.Client{}

	router := NewHandlerRouter(clientMap, emitter, store)

	msg := &api.Message{
		Room: 1,
		Type: "temperature_target",
		Data: "22.5",
	}

	ctx := context.Background()
	// Should not panic, just log error
	router.Handle(ctx, "nonexistent", msg)
}

func TestHandlerRouter_Handle_NoStoreID(t *testing.T) {
	client := mock.NewMockClient()
	emitter := &mockEmitter{}
	store := store.NewInMemoryStore()
	deviceName := "device1"

	// Don't set ID in store
	clientMap := map[string]transport.Client{
		deviceName: client,
	}

	router := NewHandlerRouter(clientMap, emitter, store)

	msg := &api.Message{
		Room: 1,
		Type: "temperature_target",
		Data: "22.5",
	}

	ctx := context.Background()
	// Should not panic, just log error
	router.Handle(ctx, deviceName, msg)
}

func TestHandlerRouter_Handle_UnknownMessageType(t *testing.T) {
	client := mock.NewMockClient()
	emitter := &mockEmitter{}
	store := store.NewInMemoryStore()
	deviceName := "device1"
	deviceID := "DEVICE-123"

	// Setup store with device ID
	store.SetID(deviceName, deviceID)

	clientMap := map[string]transport.Client{
		deviceName: client,
	}

	router := NewHandlerRouter(clientMap, emitter, store)

	msg := &api.Message{
		Room: 1,
		Type: "unknown_type",
		Data: "some data",
	}

	ctx := context.Background()
	// Should not panic, just log error
	router.Handle(ctx, deviceName, msg)
}

func TestHandlerRouter_Route_TemperatureTarget(t *testing.T) {
	client := mock.NewMockClient()
	emitter := &mockEmitter{}
	store := store.NewInMemoryStore()
	deviceID := "DEVICE-123"

	router := NewHandlerRouter(map[string]transport.Client{}, emitter, store)

	msg := &api.Message{
		Room: 2,
		Type: "temperature_target",
		Data: "23.0",
	}

	err := router.route(client, deviceID, msg)
	assert.NoError(t, err)

	// Verify the message was sent
	result, _ := client.Connect()
	assert.NotNil(t, result.Device.HeatAreas)

	found := false
	for _, heatArea := range *result.Device.HeatAreas {
		if *heatArea.Nr == 2 {
			assert.Equal(t, 23.0, *heatArea.TTarget)
			found = true
			break
		}
	}
	assert.True(t, found, "Heat area 2 should have updated temperature")
}

func TestHandlerRouter_Route_HeatareaMode(t *testing.T) {
	client := mock.NewMockClient()
	emitter := &mockEmitter{}
	store := store.NewInMemoryStore()
	deviceID := "DEVICE-123"

	router := NewHandlerRouter(map[string]transport.Client{}, emitter, store)

	tests := []struct {
		name         string
		data         string
		expectedMode int
	}{
		{"auto mode", "auto", 0},
		{"day mode", "day", 1},
		{"night mode", "night", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &api.Message{
				Room: 2,
				Type: "heatarea_mode",
				Data: tt.data,
			}

			err := router.route(client, deviceID, msg)
			assert.NoError(t, err)

			// Verify the message was sent
			result, _ := client.Connect()
			assert.NotNil(t, result.Device.HeatAreas)
			found := false
			for _, heatArea := range *result.Device.HeatAreas {
				if *heatArea.Nr == 2 {
					assert.Equal(t, tt.expectedMode, *heatArea.Mode)
					found = true
					break
				}
			}
			assert.True(t, found, "Heat area 2 should have updated mode")
		})
	}
}

func TestHandlerRouter_Route_UnknownType(t *testing.T) {
	client := mock.NewMockClient()
	emitter := &mockEmitter{}
	store := store.NewInMemoryStore()
	deviceID := "DEVICE-123"

	router := NewHandlerRouter(map[string]transport.Client{}, emitter, store)

	msg := &api.Message{
		Room: 1,
		Type: "unknown_type",
		Data: "data",
	}

	err := router.route(client, deviceID, msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown message type")
}
