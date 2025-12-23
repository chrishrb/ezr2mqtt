package mock

import (
	"testing"

	"github.com/chrishrb/ezr2mqtt/transport"
)

func TestNewMockClient(t *testing.T) {
	client := NewMockClient()

	if client == nil {
		t.Fatal("NewMockClient returned nil")
	}

	if client.currentMessage == nil {
		t.Fatal("currentMessage is nil")
	}
}

func TestConnect(t *testing.T) {
	client := NewMockClient()

	msg, err := client.Connect()

	if err != nil {
		t.Fatalf("Connect returned error: %v", err)
	}

	if msg == nil {
		t.Fatal("Connect returned nil message")
	}

	// Verify basic fields are populated
	if msg.Device.ID != "MOCK-12345" {
		t.Errorf("Expected Device.ID to be 'MOCK-12345', got '%s'", msg.Device.ID)
	}

	if msg.Device.Type != "EZR" {
		t.Errorf("Expected Device.Type to be 'EZR', got '%s'", msg.Device.Type)
	}

	if msg.Device.Name != "Mock Device" {
		t.Errorf("Expected Device.Name to be 'Mock Device', got '%s'", msg.Device.Name)
	}

	// Verify heat areas are populated
	if len(msg.Device.HeatAreas) != 2 {
		t.Errorf("Expected 2 heat areas, got %d", len(msg.Device.HeatAreas))
	}

	if len(msg.Device.HeatAreas) > 0 {
		if msg.Device.HeatAreas[0].Name != "Living Room" {
			t.Errorf("Expected first heat area to be 'Living Room', got '%s'", msg.Device.HeatAreas[0].Name)
		}

		if msg.Device.HeatAreas[0].TTarget != 22.0 {
			t.Errorf("Expected first heat area TTarget to be 22.0, got %f", msg.Device.HeatAreas[0].TTarget)
		}
	}
}

func TestSend_UpdateSingleField(t *testing.T) {
	client := NewMockClient()

	// Get initial state
	initial, _ := client.Connect()
	initialName := initial.Device.Name

	// Create a message with only one field updated
	updateMsg := &transport.Message{
		Device: transport.Device{
			Mode: 2, // Change mode from 1 to 2
		},
	}

	err := client.Send(updateMsg)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	// Verify the mode was updated
	if client.currentMessage.Device.Mode != 2 {
		t.Errorf("Expected Mode to be 2, got %d", client.currentMessage.Device.Mode)
	}

	// Verify other fields remain unchanged
	if client.currentMessage.Device.Name != initialName {
		t.Errorf("Expected Name to remain '%s', got '%s'", initialName, client.currentMessage.Device.Name)
	}

	if client.currentMessage.Device.ID != "MOCK-12345" {
		t.Errorf("Expected ID to remain 'MOCK-12345', got '%s'", client.currentMessage.Device.ID)
	}
}

func TestSend_UpdateNestedStruct(t *testing.T) {
	client := NewMockClient()

	// Get initial network config
	initial, _ := client.Connect()
	initialMAC := initial.Device.Network.MAC
	initialDHCP := initial.Device.Network.DHCP

	// Update only the IPv4 address
	updateMsg := &transport.Message{
		Device: transport.Device{
			Network: transport.Network{
				IPv4Actual: "192.168.2.100",
			},
		},
	}

	err := client.Send(updateMsg)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	// Verify IPv4 was updated
	if client.currentMessage.Device.Network.IPv4Actual != "192.168.2.100" {
		t.Errorf("Expected IPv4Actual to be '192.168.2.100', got '%s'", client.currentMessage.Device.Network.IPv4Actual)
	}

	// Verify other network fields remain unchanged
	if client.currentMessage.Device.Network.MAC != initialMAC {
		t.Errorf("Expected MAC to remain '%s', got '%s'", initialMAC, client.currentMessage.Device.Network.MAC)
	}

	if client.currentMessage.Device.Network.DHCP != initialDHCP {
		t.Errorf("Expected DHCP to remain %d, got %d", initialDHCP, client.currentMessage.Device.Network.DHCP)
	}
}

func TestSend_UpdateHeatAreaByNr(t *testing.T) {
	client := NewMockClient()

	// Get initial state
	initial, _ := client.Connect()
	initialLivingRoomTarget := initial.Device.HeatAreas[0].TTarget
	initialBedroomTarget := initial.Device.HeatAreas[1].TTarget

	// Update only heat area #1 (Living Room) target temperature
	updateMsg := &transport.Message{
		Device: transport.Device{
			HeatAreas: []transport.HeatArea{
				{
					Nr:      1,
					TTarget: 24.5,
				},
			},
		},
	}

	err := client.Send(updateMsg)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	// Verify heat area #1 was updated
	if len(client.currentMessage.Device.HeatAreas) < 2 {
		t.Fatalf("Expected at least 2 heat areas, got %d", len(client.currentMessage.Device.HeatAreas))
	}

	heatArea1 := findHeatAreaByNr(client.currentMessage.Device.HeatAreas, 1)
	if heatArea1 == nil {
		t.Fatal("Heat area #1 not found")
	}

	if heatArea1.TTarget != 24.5 {
		t.Errorf("Expected heat area #1 TTarget to be 24.5, got %f", heatArea1.TTarget)
	}

	// Verify heat area #1's other fields remain unchanged
	if heatArea1.Name != "Living Room" {
		t.Errorf("Expected heat area #1 Name to remain 'Living Room', got '%s'", heatArea1.Name)
	}

	// Verify heat area #2 remains completely unchanged
	heatArea2 := findHeatAreaByNr(client.currentMessage.Device.HeatAreas, 2)
	if heatArea2 == nil {
		t.Fatal("Heat area #2 not found")
	}

	if heatArea2.TTarget != initialBedroomTarget {
		t.Errorf("Expected heat area #2 TTarget to remain %f, got %f", initialBedroomTarget, heatArea2.TTarget)
	}

	if heatArea2.Name != "Bedroom" {
		t.Errorf("Expected heat area #2 Name to remain 'Bedroom', got '%s'", heatArea2.Name)
	}

	// Verify initial living room target was different
	if initialLivingRoomTarget == 24.5 {
		t.Error("Initial living room target should not be 24.5 for this test to be valid")
	}
}

func TestSend_UpdateMultipleHeatAreas(t *testing.T) {
	client := NewMockClient()

	// Update both heat areas
	updateMsg := &transport.Message{
		Device: transport.Device{
			HeatAreas: []transport.HeatArea{
				{
					Nr:      1,
					TTarget: 23.0,
				},
				{
					Nr:      2,
					TTarget: 19.5,
				},
			},
		},
	}

	err := client.Send(updateMsg)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	// Verify both heat areas were updated
	heatArea1 := findHeatAreaByNr(client.currentMessage.Device.HeatAreas, 1)
	if heatArea1 == nil {
		t.Fatal("Heat area #1 not found")
	}

	if heatArea1.TTarget != 23.0 {
		t.Errorf("Expected heat area #1 TTarget to be 23.0, got %f", heatArea1.TTarget)
	}

	heatArea2 := findHeatAreaByNr(client.currentMessage.Device.HeatAreas, 2)
	if heatArea2 == nil {
		t.Fatal("Heat area #2 not found")
	}

	if heatArea2.TTarget != 19.5 {
		t.Errorf("Expected heat area #2 TTarget to be 19.5, got %f", heatArea2.TTarget)
	}
}

func TestSend_AddNewHeatArea(t *testing.T) {
	client := NewMockClient()

	// Get initial count
	initial, _ := client.Connect()
	initialCount := len(initial.Device.HeatAreas)

	// Add a new heat area
	updateMsg := &transport.Message{
		Device: transport.Device{
			HeatAreas: []transport.HeatArea{
				{
					Nr:      3,
					Name:    "Kitchen",
					TTarget: 21.0,
				},
			},
		},
	}

	err := client.Send(updateMsg)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	// Verify new heat area was added
	if len(client.currentMessage.Device.HeatAreas) != initialCount+1 {
		t.Errorf("Expected %d heat areas, got %d", initialCount+1, len(client.currentMessage.Device.HeatAreas))
	}

	heatArea3 := findHeatAreaByNr(client.currentMessage.Device.HeatAreas, 3)
	if heatArea3 == nil {
		t.Fatal("Heat area #3 not found")
	}

	if heatArea3.Name != "Kitchen" {
		t.Errorf("Expected heat area #3 Name to be 'Kitchen', got '%s'", heatArea3.Name)
	}

	if heatArea3.TTarget != 21.0 {
		t.Errorf("Expected heat area #3 TTarget to be 21.0, got %f", heatArea3.TTarget)
	}
}

func TestSend_UpdateFloatValue(t *testing.T) {
	client := NewMockClient()

	// Update a float value
	updateMsg := &transport.Message{
		Device: transport.Device{
			AntifreezeTemp: 7.5,
		},
	}

	err := client.Send(updateMsg)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	// Verify float was updated
	if client.currentMessage.Device.AntifreezeTemp != 7.5 {
		t.Errorf("Expected AntifreezeTemp to be 7.5, got %f", client.currentMessage.Device.AntifreezeTemp)
	}
}

func TestSend_UpdateStringValue(t *testing.T) {
	client := NewMockClient()

	// Update a string value
	updateMsg := &transport.Message{
		Device: transport.Device{
			Name: "Updated Device Name",
		},
	}

	err := client.Send(updateMsg)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	// Verify string was updated
	if client.currentMessage.Device.Name != "Updated Device Name" {
		t.Errorf("Expected Name to be 'Updated Device Name', got '%s'", client.currentMessage.Device.Name)
	}
}

func TestSend_MultipleSequentialUpdates(t *testing.T) {
	client := NewMockClient()

	// First update
	updateMsg1 := &transport.Message{
		Device: transport.Device{
			Mode: 2,
		},
	}

	err := client.Send(updateMsg1)
	if err != nil {
		t.Fatalf("First Send returned error: %v", err)
	}

	if client.currentMessage.Device.Mode != 2 {
		t.Errorf("After first update, expected Mode to be 2, got %d", client.currentMessage.Device.Mode)
	}

	// Second update
	updateMsg2 := &transport.Message{
		Device: transport.Device{
			Cooling: 1,
		},
	}

	err = client.Send(updateMsg2)
	if err != nil {
		t.Fatalf("Second Send returned error: %v", err)
	}

	if client.currentMessage.Device.Cooling != 1 {
		t.Errorf("After second update, expected Cooling to be 1, got %d", client.currentMessage.Device.Cooling)
	}

	// Verify first update persisted
	if client.currentMessage.Device.Mode != 2 {
		t.Errorf("After second update, expected Mode to still be 2, got %d", client.currentMessage.Device.Mode)
	}
}

func TestSend_ZeroValuesDoNotOverwrite(t *testing.T) {
	client := NewMockClient()

	// Get initial mode
	initial, _ := client.Connect()
	initialMode := initial.Device.Mode

	// Try to send a message with zero mode (should not overwrite)
	updateMsg := &transport.Message{
		Device: transport.Device{
			Mode: 0,
			Name: "Updated Name",
		},
	}

	err := client.Send(updateMsg)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	// Verify mode was NOT updated (because 0 is a zero value)
	if client.currentMessage.Device.Mode != initialMode {
		t.Errorf("Expected Mode to remain %d, got %d", initialMode, client.currentMessage.Device.Mode)
	}

	// Verify name WAS updated
	if client.currentMessage.Device.Name != "Updated Name" {
		t.Errorf("Expected Name to be 'Updated Name', got '%s'", client.currentMessage.Device.Name)
	}
}

func TestSend_NilMessage(t *testing.T) {
	client := NewMockClient()

	// Get initial state
	initial, _ := client.Connect()
	initialID := initial.Device.ID

	// Send nil message (should not panic)
	err := client.Send(nil)
	if err != nil {
		t.Fatalf("Send with nil message returned error: %v", err)
	}

	// Verify state unchanged
	if client.currentMessage.Device.ID != initialID {
		t.Error("State should not change when sending nil message")
	}
}

func TestSend_UpdateHeatCtrlByNr(t *testing.T) {
	client := NewMockClient()

	// Update heat controller #1
	updateMsg := &transport.Message{
		Device: transport.Device{
			HeatCtrls: []transport.HeatCtrl{
				{
					Nr:    1,
					Actor: 75,
				},
			},
		},
	}

	err := client.Send(updateMsg)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	// Verify heat controller #1 was updated
	heatCtrl1 := findHeatCtrlByNr(client.currentMessage.Device.HeatCtrls, 1)
	if heatCtrl1 == nil {
		t.Fatal("Heat controller #1 not found")
	}

	if heatCtrl1.Actor != 75 {
		t.Errorf("Expected heat controller #1 Actor to be 75, got %d", heatCtrl1.Actor)
	}

	// Verify other fields remain unchanged
	if heatCtrl1.HeatAreaNr != 1 {
		t.Errorf("Expected heat controller #1 HeatAreaNr to remain 1, got %d", heatCtrl1.HeatAreaNr)
	}

	// Verify heat controller #2 remains unchanged
	heatCtrl2 := findHeatCtrlByNr(client.currentMessage.Device.HeatCtrls, 2)
	if heatCtrl2 == nil {
		t.Fatal("Heat controller #2 not found")
	}

	if heatCtrl2.Actor != 30 {
		t.Errorf("Expected heat controller #2 Actor to remain 30, got %d", heatCtrl2.Actor)
	}
}

func TestSend_UpdateVacationStruct(t *testing.T) {
	client := NewMockClient()

	// Update vacation state
	updateMsg := &transport.Message{
		Device: transport.Device{
			Vacation: transport.Vacation{
				State:     1,
				StartDate: "24.12.2025",
			},
		},
	}

	err := client.Send(updateMsg)
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	// Verify vacation state was updated
	if client.currentMessage.Device.Vacation.State != 1 {
		t.Errorf("Expected Vacation.State to be 1, got %d", client.currentMessage.Device.Vacation.State)
	}

	if client.currentMessage.Device.Vacation.StartDate != "24.12.2025" {
		t.Errorf("Expected Vacation.StartDate to be '24.12.2025', got '%s'", client.currentMessage.Device.Vacation.StartDate)
	}

	// Verify other vacation fields remain unchanged
	if client.currentMessage.Device.Vacation.StartTime != "00:00" {
		t.Errorf("Expected Vacation.StartTime to remain '00:00', got '%s'", client.currentMessage.Device.Vacation.StartTime)
	}
}

// Helper functions

func findHeatAreaByNr(heatAreas []transport.HeatArea, nr int) *transport.HeatArea {
	for i := range heatAreas {
		if heatAreas[i].Nr == nr {
			return &heatAreas[i]
		}
	}
	return nil
}

func findHeatCtrlByNr(heatCtrls []transport.HeatCtrl, nr int) *transport.HeatCtrl {
	for i := range heatCtrls {
		if heatCtrls[i].Nr == nr {
			return &heatCtrls[i]
		}
	}
	return nil
}
