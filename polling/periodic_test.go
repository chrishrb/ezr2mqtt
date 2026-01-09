package polling

import (
	"context"
	"testing"
	"time"

	"github.com/chrishrb/ezr2mqtt/store"
	"github.com/chrishrb/ezr2mqtt/transport/mock"
	"github.com/stretchr/testify/assert"
)

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
