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
