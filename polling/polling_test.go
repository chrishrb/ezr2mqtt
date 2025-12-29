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

type mockEmitter struct {
	emittedMessages      []*api.Message
	emittedIDs           []string
	emittedHADiscoveries []*api.HADiscovery
	emittedHAComponents  []api.HAComponent
}

func (m *mockEmitter) Emit(ctx context.Context, id string, message *api.Message) error {
	m.emittedIDs = append(m.emittedIDs, id)
	m.emittedMessages = append(m.emittedMessages, message)
	return nil
}

func (m *mockEmitter) EmitHADiscovery(ctx context.Context, component api.HAComponent, message *api.HADiscovery) error {
	m.emittedHAComponents = append(m.emittedHAComponents, component)
	m.emittedHADiscoveries = append(m.emittedHADiscoveries, message)
	return nil
}

func (m *mockEmitter) GetPrefix() string {
	return "ezr"
}

func TestNewPoller(t *testing.T) {
	client := mock.NewMockClient()
	emitter := &mockEmitter{}
	store := store.NewInMemoryStore()
	runEvery := 5 * time.Second

	poller := NewPoller("device1", client, emitter, runEvery, store)

	assert.NotNil(t, poller)
	assert.Equal(t, "device1", poller.name)
	assert.Equal(t, client, poller.client)
	assert.Equal(t, runEvery, poller.runEvery)
	assert.Equal(t, store, poller.store)
	assert.Len(t, emitter.emittedMessages, 0) // Should not be called during construction
}

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

func TestPoller_PollPeriodic_EmitsMessages(t *testing.T) {
	client := mock.NewMockClient()
	store := store.NewInMemoryStore()
	deviceName := "test-device"

	emitter := &mockEmitter{}

	// Use a very short polling interval for testing
	poller := NewPoller(deviceName, client, emitter, 50*time.Millisecond, store)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	poller.pollPeriodic(ctx)

	// Should have emitted messages for at least one poll cycle
	// Each cycle emits 2 messages per heat area (target and actual)
	// Mock client has 2 heat areas, so 4 messages per cycle
	assert.GreaterOrEqual(t, len(emitter.emittedMessages), 4)

	// Verify message types and structure
	targetFound := false
	actualFound := false
	heatareaModeFound := false

	for i, msg := range emitter.emittedMessages {
		assert.Equal(t, deviceName, emitter.emittedIDs[i])
		assert.Contains(t, []string{"temperature_target", "temperature_actual", "heatarea_mode"}, msg.Type)

		if msg.Type == "temperature_target" {
			targetFound = true
			assert.IsType(t, "19.00", msg.Data)
		}
		if msg.Type == "temperature_actual" {
			actualFound = true
			assert.IsType(t, "19.00", msg.Data)
		}
		if msg.Type == "heatarea_mode" {
			heatareaModeFound = true
			assert.IsType(t, "auto", msg.Data)
		}
	}

	assert.True(t, targetFound, "Should emit temperature_target messages")
	assert.True(t, actualFound, "Should emit temperature_actual messages")
	assert.True(t, heatareaModeFound, "Should emit heatarea_mode messages")
}

func TestPoller_PollPeriodic_ContextCancellation(t *testing.T) {
	client := mock.NewMockClient()
	store := store.NewInMemoryStore()

	emitter := &mockEmitter{}

	poller := NewPoller("device1", client, emitter, 100*time.Millisecond, store)

	ctx, cancel := context.WithCancel(context.Background())

	// Start polling in a goroutine
	done := make(chan bool)
	go func() {
		poller.pollPeriodic(ctx)
		done <- true
	}()

	// Cancel context after a short time
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for pollPeriodic to finish (with timeout)
	select {
	case <-done:
		// Test passes - pollPeriodic returned
	case <-time.After(1 * time.Second):
		t.Fatal("pollPeriodic did not respect context cancellation")
	}
}

func TestPoller_Run_StartsPolling(t *testing.T) {
	client := mock.NewMockClient()
	store := store.NewInMemoryStore()

	emitter := &mockEmitter{}

	poller := NewPoller("device1", client, emitter, 50*time.Millisecond, store)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	poller.Run(ctx)

	// Wait for polling to occur
	time.Sleep(250 * time.Millisecond)

	// Should have emitted HA discovery messages from pollOnce
	// and some periodic temperature/mode messages
	assert.Greater(t, len(emitter.emittedMessages), 0, "Should have periodic messages")
	assert.Greater(t, len(emitter.emittedHADiscoveries), 0, "Should have HA discovery messages from pollOnce")
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

func TestPoller_PollPeriodic_EmitsCorrectData(t *testing.T) {
	client := mock.NewMockClient()
	store := store.NewInMemoryStore()

	emitter := &mockEmitter{}

	poller := NewPoller("device1", client, emitter, 50*time.Millisecond, store)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	poller.pollPeriodic(ctx)

	// Verify we got messages for both heat areas
	room1Target := false
	room1Actual := false
	room2Target := false
	room2Actual := false

	for _, msg := range emitter.emittedMessages {
		if msg.Room == 1 && msg.Type == "temperature_target" {
			room1Target = true
			assert.Equal(t, "22.00", msg.Data)
		}
		if msg.Room == 1 && msg.Type == "temperature_actual" {
			room1Actual = true
			assert.Equal(t, "22.50", msg.Data)
		}
		if msg.Room == 2 && msg.Type == "temperature_target" {
			room2Target = true
			assert.Equal(t, "20.00", msg.Data)
		}
		if msg.Room == 2 && msg.Type == "temperature_actual" {
			room2Actual = true
			assert.Equal(t, "19.50", msg.Data)
		}
	}

	assert.True(t, room1Target, "Should emit room 1 target temperature")
	assert.True(t, room1Actual, "Should emit room 1 actual temperature")
	assert.True(t, room2Target, "Should emit room 2 target temperature")
	assert.True(t, room2Actual, "Should emit room 2 actual temperature")
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
	assert.Contains(t, targetDiscovery.Name, "Temperature Target")
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
	assert.Contains(t, actualDiscovery.Name, "Temperature Actual")
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
	assert.Contains(t, modeDiscovery.Name, "Heatarea Mode")
	assert.Contains(t, modeDiscovery.UniqueID, "heatarea_mode")
	assert.Contains(t, modeDiscovery.StateTopic, "state/heatarea_mode")
	assert.Contains(t, modeDiscovery.CommandTopic, "set/heatarea_mode")
	assert.Equal(t, []string{"auto", "day", "night"}, modeDiscovery.Options)
}
