package mock

import (
	"testing"

	"github.com/chrishrb/ezr2mqtt/transport"
	"github.com/stretchr/testify/assert"
)

func TestNewMockClient(t *testing.T) {
	client := NewMockClient()
	assert.NotNil(t, client, "NewMockClient returned nil")
	assert.NotNil(t, client.currentMessage, "NewMockClient did not initialize currentMessage")
}

func TestConnect(t *testing.T) {
	client := NewMockClient()

	msg, err := client.Connect()
	assert.NoError(t, err, "Connect returned error")
	assert.NotNil(t, msg, "Connect returned nil message")

	// Verify mock data is populated
	assert.NotNil(t, msg.Device.ID, "Expected Device.ID to be non-nil")
	assert.Equal(t, "MOCK-12345", *msg.Device.ID, "Expected Device.ID to be MOCK-12345")
	assert.NotNil(t, msg.Device.Type, "Expected Device.Type to be non-nil")
	assert.Equal(t, "EZR", *msg.Device.Type, "Expected Device.Type to be EZR")
	assert.NotNil(t, msg.Device.Network, "Expected Network to be populated")
	assert.NotNil(t, msg.Device.Vacation, "Expected Vacation to be populated")
	assert.NotNil(t, msg.Device.HeatAreas, "Expected HeatAreas to be non-nil")
	assert.NotEmpty(t, *msg.Device.HeatAreas, "Expected HeatAreas to be populated")
}

func TestSend_UpdateSimpleField(t *testing.T) {
	client := NewMockClient()
	_, err := client.Connect()
	assert.NoError(t, err, "Connect failed")

	// Store original values
	originalID := *client.currentMessage.Device.ID
	originalType := *client.currentMessage.Device.Type

	// Send partial update - only change Name
	newName := "Updated Name"
	partialMsg := &transport.Message{
		Device: transport.Device{
			Name: &newName,
		},
	}

	err = client.Send(partialMsg)
	assert.NoError(t, err, "Send returned error")

	// Verify only Name was updated
	assert.NotNil(t, client.currentMessage.Device.Name, "Name should not be nil")
	assert.Equal(t, "Updated Name", *client.currentMessage.Device.Name, "Name was not updated correctly")

	// Verify other fields remain unchanged
	assert.Equal(t, originalID, *client.currentMessage.Device.ID, "ID should not have changed")
	assert.Equal(t, originalType, *client.currentMessage.Device.Type, "Type should not have changed")
}

func TestSend_UpdateNestedVacation(t *testing.T) {
	client := NewMockClient()
	_, err := client.Connect()
	assert.NoError(t, err, "Connect failed")

	// Store original values
	originalStartDate := *client.currentMessage.Device.Vacation.StartDate
	originalEndDate := *client.currentMessage.Device.Vacation.EndDate

	// Send partial update - only change Vacation.State
	newState := 1
	partialMsg := &transport.Message{
		Device: transport.Device{
			Vacation: &transport.Vacation{
				State: &newState,
			},
		},
	}

	err = client.Send(partialMsg)
	assert.NoError(t, err, "Send returned error")

	// Verify only Vacation.State was updated
	assert.NotNil(t, client.currentMessage.Device.Vacation.State, "Vacation.State should not be nil")
	assert.Equal(t, 1, *client.currentMessage.Device.Vacation.State, "Vacation.State was not updated correctly")

	// Verify other Vacation fields remain unchanged
	assert.Equal(t, originalStartDate, *client.currentMessage.Device.Vacation.StartDate, "Vacation.StartDate should not have changed")
	assert.Equal(t, originalEndDate, *client.currentMessage.Device.Vacation.EndDate, "Vacation.EndDate should not have changed")
}

func TestSend_UpdateNestedNetwork(t *testing.T) {
	client := NewMockClient()
	_, err := client.Connect()
	assert.NoError(t, err, "Connect failed")

	// Store original values
	originalMAC := *client.currentMessage.Device.Network.MAC
	originalGateway := *client.currentMessage.Device.Network.Gateway

	// Send partial update - only change Network.IPv4Actual
	newIP := "192.168.1.200"
	partialMsg := &transport.Message{
		Device: transport.Device{
			Network: &transport.Network{
				IPv4Actual: &newIP,
			},
		},
	}

	err = client.Send(partialMsg)
	assert.NoError(t, err, "Send returned error")

	// Verify only Network.IPv4Actual was updated
	assert.NotNil(t, client.currentMessage.Device.Network.IPv4Actual, "Network.IPv4Actual should not be nil")
	assert.Equal(t, "192.168.1.200", *client.currentMessage.Device.Network.IPv4Actual, "Network.IPv4Actual was not updated correctly")

	// Verify other Network fields remain unchanged
	assert.Equal(t, originalMAC, *client.currentMessage.Device.Network.MAC, "Network.MAC should not have changed")
	assert.Equal(t, originalGateway, *client.currentMessage.Device.Network.Gateway, "Network.Gateway should not have changed")
}

func TestSend_UpdateNestedCloud(t *testing.T) {
	client := NewMockClient()
	_, err := client.Connect()
	assert.NoError(t, err, "Connect failed")

	// Cloud should be nil initially
	assert.Nil(t, client.currentMessage.Device.Cloud, "Cloud should be nil initially")

	// Send partial update - create Cloud with UserID
	userID := "test-user"
	partialMsg := &transport.Message{
		Device: transport.Device{
			Cloud: &transport.Cloud{
				UserID: &userID,
			},
		},
	}

	err = client.Send(partialMsg)
	assert.NoError(t, err, "Send returned error")

	// Verify Cloud.UserID was set
	assert.NotNil(t, client.currentMessage.Device.Cloud, "Cloud should not be nil after Send")
	assert.NotNil(t, client.currentMessage.Device.Cloud.UserID, "Cloud.UserID should not be nil")
	assert.Equal(t, "test-user", *client.currentMessage.Device.Cloud.UserID, "Cloud.UserID was not updated correctly")

	// Now update another field in Cloud
	port := 8080
	partialMsg2 := &transport.Message{
		Device: transport.Device{
			Cloud: &transport.Cloud{
				M2MServerPort: &port,
			},
		},
	}

	err = client.Send(partialMsg2)
	assert.NoError(t, err, "Send returned error")

	// Verify Cloud.M2MServerPort was set and UserID is still there
	assert.Equal(t, 8080, *client.currentMessage.Device.Cloud.M2MServerPort, "Cloud.M2MServerPort was not updated correctly")
	assert.Equal(t, "test-user", *client.currentMessage.Device.Cloud.UserID, "Cloud.UserID should still be test-user")
}

func TestSend_MultipleUpdates(t *testing.T) {
	client := NewMockClient()
	_, err := client.Connect()
	assert.NoError(t, err, "Connect failed")

	// First update
	name1 := "Name 1"
	msg1 := &transport.Message{
		Device: transport.Device{
			Name: &name1,
		},
	}
	err = client.Send(msg1)

	assert.NoError(t, err, "Send returned error")
	assert.Equal(t, "Name 1", *client.currentMessage.Device.Name, "First update failed")

	// Second update
	name2 := "Name 2"
	mode := 2
	msg2 := &transport.Message{
		Device: transport.Device{
			Name: &name2,
			Mode: &mode,
		},
	}
	err = client.Send(msg2)
	assert.NoError(t, err, "Send returned error")

	assert.Equal(t, "Name 2", *client.currentMessage.Device.Name, "Second update failed for Name")
	assert.Equal(t, 2, *client.currentMessage.Device.Mode, "Second update failed for Mode")
}

func TestSend_WithNilCurrentMessage(t *testing.T) {
	client := &MockClient{
		currentMessage: nil,
	}

	name := "Test Name"
	msg := &transport.Message{
		Device: transport.Device{
			Name: &name,
		},
	}

	err := client.Send(msg)
	assert.NoError(t, err, "Send returned error")

	assert.NotNil(t, client.currentMessage, "currentMessage should be initialized")
	assert.NotNil(t, client.currentMessage.Device.Name, "Name should not be nil")
	assert.Equal(t, "Test Name", *client.currentMessage.Device.Name, "Name was not set correctly")
}

func TestSend_UpdateHeatAreas(t *testing.T) {
	client := NewMockClient()
	_, err := client.Connect()
	assert.NoError(t, err, "Connect failed")

	// Store original heat area
	originalLen := len(*client.currentMessage.Device.HeatAreas)
	assert.Equal(t, 2, originalLen, "Expected 2 heat area initially")
	originalName := *(*client.currentMessage.Device.HeatAreas)[0].Name

	// Update existing heat area (Nr=1) with only TTarget field
	newTarget := 20.5
	updatedHeatAreas := []transport.HeatArea{
		{
			Nr:      intPtr(1),
			TTarget: &newTarget,
		},
	}

	msg := &transport.Message{
		Device: transport.Device{
			HeatAreas: &updatedHeatAreas,
		},
	}

	err = client.Send(msg)
	assert.NoError(t, err, "Send returned error")

	// Verify heat area count is still 1 (merged, not replaced)
	assert.Equal(t, 2, len(*client.currentMessage.Device.HeatAreas), "Expected 2 heat area after update")

	// Verify TTarget was updated
	assert.Equal(t, newTarget, *(*client.currentMessage.Device.HeatAreas)[0].TTarget, "TTarget should be updated to 20.5")

	// Verify original Name is preserved
	assert.Equal(t, originalName, *(*client.currentMessage.Device.HeatAreas)[0].Name, "Original Name should be preserved")
}

func TestSend_AddNewHeatArea(t *testing.T) {
	client := NewMockClient()
	_, err := client.Connect()
	assert.NoError(t, err, "Connect failed")

	// Add a new heat area with Nr=2
	newHeatAreas := []transport.HeatArea{
		{
			Nr:      intPtr(2),
			Name:    strPtr("Bedroom"),
			TTarget: floatPtr(20.0),
		},
	}

	msg := &transport.Message{
		Device: transport.Device{
			HeatAreas: &newHeatAreas,
		},
	}

	err = client.Send(msg)
	assert.NoError(t, err, "Send returned error")

	// Verify we now have 2 heat areas
	assert.Equal(t, 2, len(*client.currentMessage.Device.HeatAreas), "Expected 2 heat areas")

	// Verify the new area was added
	found := false
	for _, area := range *client.currentMessage.Device.HeatAreas {
		if area.Nr != nil && *area.Nr == 2 {
			found = true
			assert.Equal(t, "Bedroom", *area.Name, "New heat area name is incorrect")
			break
		}
	}
	assert.True(t, found, "New heat area with Nr=2 was not added")
}

func TestSend_MergeMultipleHeatAreas(t *testing.T) {
	client := NewMockClient()
	_, err := client.Connect()
	assert.NoError(t, err, "Connect failed")

	// Send update with both existing and new heat areas
	newHeatAreas := []transport.HeatArea{
		{
			Nr:      intPtr(1), // Update existing
			TTarget: floatPtr(24.0),
		},
		{
			Nr:   intPtr(2), // Add new
			Name: strPtr("Kitchen"),
		},
	}

	msg := &transport.Message{
		Device: transport.Device{
			HeatAreas: &newHeatAreas,
		},
	}

	err = client.Send(msg)
	assert.NoError(t, err, "Send returned error")

	// Verify we have 2 heat areas
	assert.Equal(t, 2, len(*client.currentMessage.Device.HeatAreas), "Expected 2 heat areas")

	// Verify first area was updated, not replaced
	for _, area := range *client.currentMessage.Device.HeatAreas {
		if area.Nr != nil && *area.Nr == 1 {
			assert.Equal(t, 24.0, *area.TTarget, "HeatArea Nr=1 TTarget was not updated")
			assert.Equal(t, "Living Room", *area.Name, "HeatArea Nr=1 original Name was lost")
		}
		if area.Nr != nil && *area.Nr == 2 {
			assert.Equal(t, "Kitchen", *area.Name, "HeatArea Nr=2 was not added correctly")
		}
	}
}

func TestSend_ComplexScenario(t *testing.T) {
	client := NewMockClient()
	_, err := client.Connect()
	assert.NoError(t, err, "Connect failed")

	// Update multiple fields across different levels
	newName := "Complex Device"
	newMode := 3
	newVacState := 1
	newIP := "10.0.0.50"

	msg := &transport.Message{
		Device: transport.Device{
			Name: &newName,
			Mode: &newMode,
			Vacation: &transport.Vacation{
				State: &newVacState,
			},
			Network: &transport.Network{
				IPv4Actual: &newIP,
			},
		},
	}

	err = client.Send(msg)
	assert.NoError(t, err, "Send returned error")

	// Verify all updates
	assert.Equal(t, "Complex Device", *client.currentMessage.Device.Name, "Name was not updated")
	assert.Equal(t, 3, *client.currentMessage.Device.Mode, "Mode was not updated")
	assert.Equal(t, 1, *client.currentMessage.Device.Vacation.State, "Vacation.State was not updated")
	assert.Equal(t, "10.0.0.50", *client.currentMessage.Device.Network.IPv4Actual, "Network.IPv4Actual was not updated")

	// Verify original values are preserved
	assert.Equal(t, "MOCK-12345", *client.currentMessage.Device.ID, "Original ID was changed")
	assert.Equal(t, "00:11:22:33:44:55", *client.currentMessage.Device.Network.MAC, "Original MAC address was changed")
	assert.Equal(t, "2025-01-01", *client.currentMessage.Device.Vacation.StartDate, "Original Vacation.StartDate was changed")
}

// Helper functions for tests
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

func TestSend_UpdateNestedPumpOutput(t *testing.T) {
	client := NewMockClient()
	_, err := client.Connect()
	assert.NoError(t, err, "Connect failed")

	// PumpOutput should be nil initially
	assert.Nil(t, client.currentMessage.Device.PumpOutput, "PumpOutput should be nil initially")

	// Set one field
	leadTime := 30
	msg1 := &transport.Message{
		Device: transport.Device{
			PumpOutput: &transport.PumpOutput{
				LeadTime: &leadTime,
			},
		},
	}

	err = client.Send(msg1)
	assert.NoError(t, err, "Send returned error")

	assert.NotNil(t, client.currentMessage.Device.PumpOutput, "PumpOutput should not be nil after Send")
	assert.Equal(t, 30, *client.currentMessage.Device.PumpOutput.LeadTime, "LeadTime was not set correctly")

	// Update another field, first should be preserved
	opMode := 1
	msg2 := &transport.Message{
		Device: transport.Device{
			PumpOutput: &transport.PumpOutput{
				OperationMode: &opMode,
			},
		},
	}

	err = client.Send(msg2)
	assert.NoError(t, err, "Send returned error")

	assert.Equal(t, 30, *client.currentMessage.Device.PumpOutput.LeadTime, "LeadTime should still be 30")
	assert.Equal(t, 1, *client.currentMessage.Device.PumpOutput.OperationMode, "OperationMode was not set correctly")
}
