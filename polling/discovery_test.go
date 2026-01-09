package polling

import (
	"context"
	"testing"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/store"
	"github.com/chrishrb/ezr2mqtt/transport/mock"
	"github.com/stretchr/testify/assert"
)

func TestPoller_PollOnce_Success(t *testing.T) {
	client := mock.NewMockClient()
	store := store.NewInMemoryStore()
	deviceName := "test-device"

	emitter := &mockEmitter{}

	poller := NewPoller(deviceName, client, emitter, 1*time.Hour, store)

	ctx := context.Background()
	poller.pollOnce(ctx)

	// Give it a moment to complete
	time.Sleep(100 * time.Millisecond)

	// Verify device ID was stored
	id := store.GetID(deviceName)
	assert.NotNil(t, id)
	assert.Equal(t, "MOCK-12345", *id)

	// Verify HA discovery messages were emitted (not regular messages)
	// Mock client has 2 heat areas, each emits 3 HA discovery messages
	assert.Len(t, emitter.emittedHADiscoveries, 6)
	assert.Len(t, emitter.emittedMessages, 0) // pollOnce doesn't emit regular messages
}

func TestPoller_PollOnce_StoresCorrectDeviceID(t *testing.T) {
	client := mock.NewMockClient()
	store := store.NewInMemoryStore()
	deviceName := "my-device"

	emitter := &mockEmitter{}

	poller := NewPoller(deviceName, client, emitter, 1*time.Hour, store)

	ctx := context.Background()
	poller.pollOnce(ctx)

	time.Sleep(50 * time.Millisecond)

	// Verify the device ID was stored correctly
	id := store.GetID(deviceName)
	assert.NotNil(t, id)
	assert.Equal(t, "MOCK-12345", *id)
}

func TestPoller_PollOnce_EmitsHADiscovery(t *testing.T) {
	client := mock.NewMockClient()
	store := store.NewInMemoryStore()
	deviceName := "test-device"

	emitter := &mockEmitter{}

	poller := NewPoller(deviceName, client, emitter, 1*time.Hour, store)

	ctx := context.Background()
	poller.pollOnce(ctx)

	// Give it a moment to complete
	time.Sleep(100 * time.Millisecond)

	// Mock client has 2 heat areas, each should emit 3 HA discovery messages
	// (temperature_target as number, temperature_actual as sensor, heatarea_mode as select)
	assert.Len(t, emitter.emittedHADiscoveries, 6)
	assert.Len(t, emitter.emittedHAComponents, 6)

	// Verify the component types are correct
	componentCounts := map[api.HAComponent]int{}
	for _, component := range emitter.emittedHAComponents {
		componentCounts[component]++
	}
	assert.Equal(t, 2, componentCounts[api.HAComponentNumber], "Should have 2 number components (target temp for each room)")
	assert.Equal(t, 2, componentCounts[api.HAComponentSensor], "Should have 2 sensor components (actual temp for each room)")
	assert.Equal(t, 2, componentCounts[api.HAComponentSelect], "Should have 2 select components (mode for each room)")

	// Verify discovery message content for temperature target (number component)
	var targetDiscovery *api.HADiscovery
	for i, component := range emitter.emittedHAComponents {
		if component == api.HAComponentNumber {
			targetDiscovery = emitter.emittedHADiscoveries[i]
			break
		}
	}
	assert.NotNil(t, targetDiscovery)
	assert.Contains(t, targetDiscovery.Name, "temperature_target")
	assert.Contains(t, targetDiscovery.UniqueID, "test-device")
	assert.Contains(t, targetDiscovery.UniqueID, "temperature_target")
	assert.Contains(t, targetDiscovery.StateTopic, "ezr/test-device")
	assert.Contains(t, targetDiscovery.StateTopic, "state/temperature_target")
	assert.Equal(t, "°C", targetDiscovery.UnitOfMeasurement)
	assert.Equal(t, "temperature", targetDiscovery.DeviceClass)
	assert.Equal(t, "measurement", targetDiscovery.StateClass)
	assert.Contains(t, targetDiscovery.CommandTopic, "set/temperature_target")
	assert.Equal(t, "slider", targetDiscovery.Mode)
	assert.NotNil(t, targetDiscovery.Device)
	assert.Equal(t, "MOCK-12345", targetDiscovery.Device.Identifiers[0])
	assert.Equal(t, "Mock Device", targetDiscovery.Device.Name)

	// Verify discovery message content for temperature actual (sensor component)
	var actualDiscovery *api.HADiscovery
	for i, component := range emitter.emittedHAComponents {
		if component == api.HAComponentSensor {
			actualDiscovery = emitter.emittedHADiscoveries[i]
			break
		}
	}
	assert.NotNil(t, actualDiscovery)
	assert.Contains(t, actualDiscovery.Name, "temperature_actual")
	assert.Contains(t, actualDiscovery.UniqueID, "temperature_actual")
	assert.Contains(t, actualDiscovery.StateTopic, "state/temperature_actual")
	assert.Equal(t, "°C", actualDiscovery.UnitOfMeasurement)
	assert.Equal(t, "temperature", actualDiscovery.DeviceClass)
	assert.Equal(t, "measurement", actualDiscovery.StateClass)

	// Verify discovery message content for heatarea mode (select component)
	var modeDiscovery *api.HADiscovery
	for i, component := range emitter.emittedHAComponents {
		if component == api.HAComponentSelect {
			modeDiscovery = emitter.emittedHADiscoveries[i]
			break
		}
	}
	assert.NotNil(t, modeDiscovery)
	assert.Contains(t, modeDiscovery.Name, "heatarea_mode")
	assert.Contains(t, modeDiscovery.UniqueID, "heatarea_mode")
	assert.Contains(t, modeDiscovery.StateTopic, "state/heatarea_mode")
	assert.Contains(t, modeDiscovery.CommandTopic, "set/heatarea_mode")
	assert.Equal(t, []string{"auto", "day", "night"}, modeDiscovery.Options)
}
