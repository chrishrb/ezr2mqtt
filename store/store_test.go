package store

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryStore(t *testing.T) {
	store := NewInMemoryStore()

	assert.NotNil(t, store)
	assert.NotNil(t, store.ids)
	assert.Len(t, store.ids, 0)
}

func TestInMemoryStore_SetID(t *testing.T) {
	store := NewInMemoryStore()

	store.SetID("device1", "ID-123")

	// Verify the ID was stored
	id := store.GetID("device1")
	assert.NotNil(t, id)
	assert.Equal(t, "ID-123", *id)
}

func TestInMemoryStore_GetID_NotFound(t *testing.T) {
	store := NewInMemoryStore()

	id := store.GetID("nonexistent")
	assert.Nil(t, id)
}

func TestInMemoryStore_GetID_Found(t *testing.T) {
	store := NewInMemoryStore()
	store.SetID("device1", "ID-123")

	id := store.GetID("device1")
	assert.NotNil(t, id)
	assert.Equal(t, "ID-123", *id)
}

func TestInMemoryStore_SetID_Update(t *testing.T) {
	store := NewInMemoryStore()

	// Set initial ID
	store.SetID("device1", "ID-123")
	id := store.GetID("device1")
	assert.Equal(t, "ID-123", *id)

	// Update ID
	store.SetID("device1", "ID-456")
	id = store.GetID("device1")
	assert.NotNil(t, id)
	assert.Equal(t, "ID-456", *id)
}

func TestInMemoryStore_MultipleDevices(t *testing.T) {
	store := NewInMemoryStore()

	// Store multiple devices
	store.SetID("device1", "ID-123")
	store.SetID("device2", "ID-456")
	store.SetID("device3", "ID-789")

	// Verify all devices are stored correctly
	id1 := store.GetID("device1")
	assert.NotNil(t, id1)
	assert.Equal(t, "ID-123", *id1)

	id2 := store.GetID("device2")
	assert.NotNil(t, id2)
	assert.Equal(t, "ID-456", *id2)

	id3 := store.GetID("device3")
	assert.NotNil(t, id3)
	assert.Equal(t, "ID-789", *id3)

	// Verify non-existent device returns nil
	id4 := store.GetID("device4")
	assert.Nil(t, id4)
}

func TestInMemoryStore_ThreadSafety(t *testing.T) {
	store := NewInMemoryStore()
	var wg sync.WaitGroup

	// Number of concurrent operations
	n := 100

	// Concurrent writes
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			deviceName := "device"
			deviceID := "ID"
			store.SetID(deviceName, deviceID)
		}(i)
	}

	// Concurrent reads
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			deviceName := "device"
			_ = store.GetID(deviceName)
		}(i)
	}

	wg.Wait()

	// Verify the store is still functional
	id := store.GetID("device")
	assert.NotNil(t, id)
	assert.Equal(t, "ID", *id)
}

func TestInMemoryStore_EmptyName(t *testing.T) {
	store := NewInMemoryStore()

	// Test with empty device name
	store.SetID("", "ID-123")
	id := store.GetID("")
	assert.NotNil(t, id)
	assert.Equal(t, "ID-123", *id)
}

func TestInMemoryStore_EmptyID(t *testing.T) {
	store := NewInMemoryStore()

	// Test with empty device ID
	store.SetID("device1", "")
	id := store.GetID("device1")
	assert.NotNil(t, id)
	assert.Equal(t, "", *id)
}

func TestInMemoryStore_ConcurrentReadWrite(t *testing.T) {
	store := NewInMemoryStore()
	var wg sync.WaitGroup

	deviceCount := 10
	opsPerDevice := 10

	// Create multiple devices concurrently
	for i := range deviceCount {
		for j := range opsPerDevice {
			wg.Add(1)
			go func(deviceNum, opNum int) {
				defer wg.Done()
				deviceName := "device" + string(rune('0'+deviceNum))
				deviceID := "ID-" + string(rune('0'+opNum))
				store.SetID(deviceName, deviceID)
			}(i, j)
		}
	}

	// Read concurrently while writing
	for i := range deviceCount {
		for range opsPerDevice {
			wg.Add(1)
			go func(deviceNum int) {
				defer wg.Done()
				deviceName := "device" + string(rune('0'+deviceNum))
				_ = store.GetID(deviceName)
			}(i)
		}
	}

	wg.Wait()

	// No assertion needed - the test passes if there's no race condition
}

func TestInMemoryStore_Interface(t *testing.T) {
	// Verify that InMemoryStore implements the Store interface
	var _ Store = (*InMemoryStore)(nil)
}
