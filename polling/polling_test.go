package polling

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/store"
	"github.com/chrishrb/ezr2mqtt/transport/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewPoller(t *testing.T) {
	client := mock.NewMockClient()
	emitterCalled := false
	emitter := api.EmitterFunc(func(ctx context.Context, id string, message *api.Message) error {
		emitterCalled = true
		return nil
	})
	store := store.NewInMemoryStore()
	runEvery := 5 * time.Second

	poller := NewPoller("device1", client, emitter, runEvery, store)

	assert.NotNil(t, poller)
	assert.Equal(t, "device1", poller.name)
	assert.Equal(t, client, poller.client)
	assert.Equal(t, runEvery, poller.runEvery)
	assert.Equal(t, store, poller.store)
	assert.False(t, emitterCalled) // Should not be called during construction
}

func TestPoller_PollOnce_Success(t *testing.T) {
	client := mock.NewMockClient()
	store := store.NewInMemoryStore()
	deviceName := "test-device"

	var emittedMessages []*api.Message
	var emittedIDs []string

	emitter := api.EmitterFunc(func(ctx context.Context, id string, message *api.Message) error {
		emittedIDs = append(emittedIDs, id)
		emittedMessages = append(emittedMessages, message)
		return nil
	})

	poller := NewPoller(deviceName, client, emitter, 1*time.Hour, store)

	ctx := context.Background()
	poller.pollOnce(ctx)

	// Give it a moment to complete
	time.Sleep(100 * time.Millisecond)

	// Verify device ID was stored
	id := store.GetID(deviceName)
	assert.NotNil(t, id)
	assert.Equal(t, "MOCK-12345", *id)

	// Verify meta message was emitted
	assert.Len(t, emittedMessages, 1)
	assert.Equal(t, deviceName, emittedIDs[0])
	assert.Equal(t, "meta", emittedMessages[0].Type)
	assert.Equal(t, 0, emittedMessages[0].Room)

	// Verify meta data structure
	var metaData api.ClimateDiscovery
	err := json.Unmarshal([]byte(emittedMessages[0].Data), &metaData)
	assert.NoError(t, err)
	assert.Equal(t, "Mock Device", metaData.Name)
	assert.Equal(t, "MOCK-12345", metaData.ID)
	assert.Equal(t, "EZR", metaData.Type)
	assert.Len(t, metaData.Rooms, 2)

	// Verify room data
	assert.Equal(t, 1, metaData.Rooms[0].ID)
	assert.Equal(t, "Living Room", metaData.Rooms[0].Name)
	assert.Equal(t, 2, metaData.Rooms[1].ID)
	assert.Equal(t, "Bedroom", metaData.Rooms[1].Name)
}

func TestPoller_PollPeriodic_EmitsMessages(t *testing.T) {
	client := mock.NewMockClient()
	store := store.NewInMemoryStore()
	deviceName := "test-device"

	var emittedMessages []*api.Message
	var emittedIDs []string

	emitter := api.EmitterFunc(func(ctx context.Context, id string, message *api.Message) error {
		emittedIDs = append(emittedIDs, id)
		emittedMessages = append(emittedMessages, message)
		return nil
	})

	// Use a very short polling interval for testing
	poller := NewPoller(deviceName, client, emitter, 50*time.Millisecond, store)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	poller.pollPeriodic(ctx)

	// Should have emitted messages for at least one poll cycle
	// Each cycle emits 2 messages per heat area (target and actual)
	// Mock client has 2 heat areas, so 4 messages per cycle
	assert.GreaterOrEqual(t, len(emittedMessages), 4)

	// Verify message types and structure
	targetFound := false
	actualFound := false
	heatareaModeFound := false

	for i, msg := range emittedMessages {
		assert.Equal(t, deviceName, emittedIDs[i])
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

	emitter := api.EmitterFunc(func(ctx context.Context, id string, message *api.Message) error {
		return nil
	})

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

	var emittedMessages []*api.Message

	emitter := api.EmitterFunc(func(ctx context.Context, id string, message *api.Message) error {
		emittedMessages = append(emittedMessages, message)
		return nil
	})

	poller := NewPoller("device1", client, emitter, 50*time.Millisecond, store)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	poller.Run(ctx)

	// Wait for polling to occur
	time.Sleep(250 * time.Millisecond)

	// Should have emitted at least the meta message from pollOnce
	// and some periodic messages
	assert.Greater(t, len(emittedMessages), 0)

	// First message should be meta
	assert.Equal(t, "meta", emittedMessages[0].Type)
}

func TestPoller_PollOnce_StoresCorrectDeviceID(t *testing.T) {
	client := mock.NewMockClient()
	store := store.NewInMemoryStore()
	deviceName := "my-device"

	emitter := api.EmitterFunc(func(ctx context.Context, id string, message *api.Message) error {
		return nil
	})

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

	var emittedMessages []*api.Message

	emitter := api.EmitterFunc(func(ctx context.Context, id string, message *api.Message) error {
		emittedMessages = append(emittedMessages, message)
		return nil
	})

	poller := NewPoller("device1", client, emitter, 50*time.Millisecond, store)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	poller.pollPeriodic(ctx)

	// Verify we got messages for both heat areas
	room1Target := false
	room1Actual := false
	room2Target := false
	room2Actual := false

	for _, msg := range emittedMessages {
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
