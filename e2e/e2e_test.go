//go:build e2e

package e2e_test

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/api/mqtt"
	"github.com/chrishrb/ezr2mqtt/handlers"
	"github.com/chrishrb/ezr2mqtt/polling"
	"github.com/chrishrb/ezr2mqtt/store"
	"github.com/chrishrb/ezr2mqtt/transport"
	"github.com/chrishrb/ezr2mqtt/transport/mock"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_SetTargetTemperatureOverMQTT(t *testing.T) {
	// Start MQTT broker
	broker, brokerURL := mqtt.NewBroker(t)
	go func() {
		err := broker.Serve()
		if err != nil {
			t.Logf("broker serve error: %v", err)
		}
	}()
	defer func() {
		err := broker.Close()
		if err != nil {
			t.Logf("broker close error: %v", err)
		}
	}()

	// Wait for broker to be ready
	time.Sleep(100 * time.Millisecond)

	const (
		deviceName = "test-device"
		mqttPrefix = "ezr"
		roomNr     = 1
	)

	// Create mock EZR client
	mockClient := mock.NewMockClient()

	// Create store
	memStore := store.NewInMemoryStore()

	// Get initial state to store device ID
	initialMsg, err := mockClient.Connect()
	require.NoError(t, err)
	memStore.SetID(deviceName, initialMsg.Device.ID)

	// Create handler router
	clients := map[string]transport.Client{
		deviceName: mockClient,
	}
	handlerRouter := handlers.NewHandlerRouter(clients, memStore)

	// Create MQTT listener
	listener := mqtt.NewListener(
		mqtt.WithMqttBrokerUrl[mqtt.Listener](brokerURL),
		mqtt.WithMqttPrefix[mqtt.Listener](mqttPrefix),
	)

	// Connect listener
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := listener.Connect(ctx, handlerRouter)
	require.NoError(t, err)
	defer func() {
		err := conn.Disconnect(context.Background())
		if err != nil {
			t.Logf("disconnect error: %v", err)
		}
	}()

	// Give listener time to subscribe
	time.Sleep(200 * time.Millisecond)

	// Create a test MQTT client to publish temperature change
	testClient, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{brokerURL},
		KeepAlive:         10,
		ConnectRetryDelay: 1 * time.Second,
		ClientConfig: paho.ClientConfig{
			ClientID: "test-client",
		},
	})
	require.NoError(t, err)
	defer func() {
		err := testClient.Disconnect(context.Background())
		if err != nil {
			t.Logf("test client disconnect error: %v", err)
		}
	}()

	err = testClient.AwaitConnection(ctx)
	require.NoError(t, err)

	// Prepare message to set target temperature
	newTargetTemp := 23.5

	// Publish temperature change to MQTT
	// Topic format: ezr/{device_id}/{room_nr}/set/temperature_target
	setTopic := fmt.Sprintf("%s/%s/%d/set/temperature_target", mqttPrefix, deviceName, roomNr)
	_, err = testClient.Publish(ctx, &paho.Publish{
		Topic:   setTopic,
		Payload: []byte(api.FormatFloat(newTargetTemp)),
	})
	require.NoError(t, err)

	// Wait for message to be processed
	time.Sleep(300 * time.Millisecond)

	// Verify the temperature was set in the mock client
	updatedMsg, err := mockClient.Connect()
	require.NoError(t, err)

	// Find the heat area and verify temperature
	var found bool
	for _, heatArea := range updatedMsg.Device.HeatAreas {
		if heatArea.Nr == roomNr {
			assert.Equal(t, newTargetTemp, heatArea.TTarget, "Target temperature should be updated")
			found = true
			break
		}
	}
	assert.True(t, found, "Heat area %d should exist", roomNr)
}

func TestE2E_TemperaturePublishedToMQTT(t *testing.T) {
	// Start MQTT broker
	broker, brokerURL := mqtt.NewBroker(t)
	go func() {
		err := broker.Serve()
		if err != nil {
			t.Logf("broker serve error: %v", err)
		}
	}()
	defer func() {
		err := broker.Close()
		if err != nil {
			t.Logf("broker close error: %v", err)
		}
	}()

	// Wait for broker to be ready
	time.Sleep(100 * time.Millisecond)

	const (
		deviceName = "test-device"
		mqttPrefix = "ezr"
		pollEvery  = 500 * time.Millisecond
	)

	// Create mock EZR client
	mockClient := mock.NewMockClient()

	// Create store
	memStore := store.NewInMemoryStore()

	// Get initial state to store device ID
	initialMsg, err := mockClient.Connect()
	require.NoError(t, err)
	deviceID := initialMsg.Device.ID
	memStore.SetID(deviceName, deviceID)

	// Create MQTT emitter
	emitter := mqtt.NewEmitter(
		mqtt.WithMqttBrokerUrl[mqtt.Emitter](brokerURL),
		mqtt.WithMqttPrefix[mqtt.Emitter](mqttPrefix),
	)

	// Create a test MQTT client to subscribe to temperature updates
	messageChan := make(chan *paho.Publish, 10)
	router := paho.NewStandardRouter()

	// Subscribe to state topics for all rooms
	// Topic format: ezr/{device_id}/+/state/+
	stateTopic := fmt.Sprintf("%s/%s/+/state/+", mqttPrefix, deviceName)
	router.RegisterHandler(stateTopic, func(p *paho.Publish) {
		messageChan <- p
	})

	testClient, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{brokerURL},
		KeepAlive:         10,
		ConnectRetryDelay: 1 * time.Second,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			_, err := manager.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{{Topic: stateTopic}},
			})
			if err != nil {
				t.Logf("subscription error: %v", err)
			}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: "test-subscriber",
			Router:   router,
		},
	})
	require.NoError(t, err)
	defer func() {
		err := testClient.Disconnect(context.Background())
		if err != nil {
			t.Logf("test client disconnect error: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = testClient.AwaitConnection(ctx)
	require.NoError(t, err)

	// Give subscriber time to set up
	time.Sleep(200 * time.Millisecond)

	// Create and start poller
	poller := polling.NewPoller(deviceName, mockClient, emitter, pollEvery, memStore)
	pollerCtx := t.Context()

	poller.Run(pollerCtx)

	// Collect messages for analysis
	receivedMessages := make(map[string]*api.Message)
	timeout := time.After(3 * time.Second)

	// Define expected message types per heat area
	expectedMessageTypes := []string{"temperature_target", "temperature_actual", "heatarea_mode"}
	numRooms := len(initialMsg.Device.HeatAreas)
	numMetaMessages := 1
	expectedMinMessages := numMetaMessages + (numRooms * len(expectedMessageTypes))

collecting:
	for len(receivedMessages) < expectedMinMessages {
		select {
		case msg := <-messageChan:
			topicParts := strings.Split(msg.Topic, "/")

			tp := topicParts[len(topicParts)-1]
			room, err := strconv.Atoi(topicParts[len(topicParts)-3])
			require.NoError(t, err)

			receivedMessages[msg.Topic] = &api.Message{
				Room: room,
				Type: tp,
				Data: string(msg.Payload),
			}
		case <-timeout:
			break collecting
		}
	}

	// Verify we received messages
	assert.GreaterOrEqual(t, len(receivedMessages), expectedMinMessages,
		"Should receive at least %d messages (meta + temperature data)", expectedMinMessages)

	// Verify messages for each room
	for _, heatArea := range initialMsg.Device.HeatAreas {
		// Verify temperature_target
		targetTopic := fmt.Sprintf("%s/%s/%d/state/temperature_target", mqttPrefix, deviceName, heatArea.Nr)
		if msg, ok := receivedMessages[targetTopic]; ok {
			assert.Equal(t, heatArea.Nr, msg.Room, "Room number should match")
			assert.Equal(t, "temperature_target", msg.Type, "Message type should be temperature_target")
			assert.Equal(t, api.FormatFloat(heatArea.TTarget), msg.Data, "Target temperature should match")
		} else {
			t.Errorf("Expected to receive temperature_target message for room %d on topic %s", heatArea.Nr, targetTopic)
		}

		// Verify temperature_actual
		actualTopic := fmt.Sprintf("%s/%s/%d/state/temperature_actual", mqttPrefix, deviceName, heatArea.Nr)
		if msg, ok := receivedMessages[actualTopic]; ok {
			assert.Equal(t, heatArea.Nr, msg.Room, "Room number should match")
			assert.Equal(t, "temperature_actual", msg.Type, "Message type should be temperature_actual")
			assert.Equal(t, api.FormatFloat(heatArea.TActual), msg.Data, "Actual temperature should match")
		} else {
			t.Errorf("Expected to receive temperature_actual message for room %d on topic %s", heatArea.Nr, actualTopic)
		}

		// Verify heatarea_mode
		modeTopic := fmt.Sprintf("%s/%s/%d/state/heatarea_mode", mqttPrefix, deviceName, heatArea.Nr)
		if msg, ok := receivedMessages[modeTopic]; ok {
			assert.Equal(t, heatArea.Nr, msg.Room, "Room number should match")
			assert.Equal(t, "heatarea_mode", msg.Type, "Message type should be heatarea_mode")
			assert.Equal(t, "day", msg.Data, "Heat area mode should match")
		} else {
			t.Errorf("Expected to receive heatarea_mode message for room %d on topic %s", heatArea.Nr, modeTopic)
		}
	}

	// Verify meta message was sent
	metaTopic := fmt.Sprintf("%s/%s/0/state/meta", mqttPrefix, deviceName)
	if msg, ok := receivedMessages[metaTopic]; ok {
		assert.Equal(t, 0, msg.Room, "Meta message should have room 0")
		assert.Equal(t, "meta", msg.Type, "Message type should be meta")
		assert.NotNil(t, msg.Data, "Meta data should not be nil")
	} else {
		t.Errorf("Expected to receive meta message on topic %s", metaTopic)
	}
}

func TestE2E_SetModeOverMQTT(t *testing.T) {
	// Start MQTT broker
	broker, brokerURL := mqtt.NewBroker(t)
	go func() {
		err := broker.Serve()
		if err != nil {
			t.Logf("broker serve error: %v", err)
		}
	}()
	defer func() {
		err := broker.Close()
		if err != nil {
			t.Logf("broker close error: %v", err)
		}
	}()

	// Wait for broker to be ready
	time.Sleep(100 * time.Millisecond)

	const (
		deviceName = "test-device"
		mqttPrefix = "ezr"
		roomNr     = 1
	)

	// Create mock EZR client
	mockClient := mock.NewMockClient()

	// Create store
	memStore := store.NewInMemoryStore()

	// Get initial state to store device ID
	initialMsg, err := mockClient.Connect()
	require.NoError(t, err)
	memStore.SetID(deviceName, initialMsg.Device.ID)

	// Get initial mode for room 1
	var initialMode int
	for _, ha := range initialMsg.Device.HeatAreas {
		if ha.Nr == roomNr {
			initialMode = ha.Mode
			break
		}
	}

	// Create handler router
	clients := map[string]transport.Client{
		deviceName: mockClient,
	}
	handlerRouter := handlers.NewHandlerRouter(clients, memStore)

	// Create MQTT listener
	listener := mqtt.NewListener(
		mqtt.WithMqttBrokerUrl[mqtt.Listener](brokerURL),
		mqtt.WithMqttPrefix[mqtt.Listener](mqttPrefix),
	)

	// Connect listener
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := listener.Connect(ctx, handlerRouter)
	require.NoError(t, err)
	defer func() {
		err := conn.Disconnect(context.Background())
		if err != nil {
			t.Logf("disconnect error: %v", err)
		}
	}()

	// Give listener time to subscribe
	time.Sleep(200 * time.Millisecond)

	// Create a test MQTT client to publish mode change
	testClient, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{brokerURL},
		KeepAlive:         10,
		ConnectRetryDelay: 1 * time.Second,
		ClientConfig: paho.ClientConfig{
			ClientID: "test-client",
		},
	})
	require.NoError(t, err)
	defer func() {
		err := testClient.Disconnect(context.Background())
		if err != nil {
			t.Logf("test client disconnect error: %v", err)
		}
	}()

	err = testClient.AwaitConnection(ctx)
	require.NoError(t, err)

	// Test each mode: auto, day, night
	testCases := []struct {
		mode         string
		expectedMode int
	}{
		{"night", 2},
		{"auto", 0},
		{"day", 1},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("set_mode_%s", tc.mode), func(t *testing.T) {
			// Publish mode change to MQTT
			// Topic format: ezr/{device_id}/{room_nr}/set/heatarea_mode
			setTopic := fmt.Sprintf("%s/%s/%d/set/heatarea_mode", mqttPrefix, deviceName, roomNr)
			_, err = testClient.Publish(ctx, &paho.Publish{
				Topic:   setTopic,
				Payload: []byte(tc.mode),
			})
			require.NoError(t, err)

			// Wait for message to be processed
			time.Sleep(300 * time.Millisecond)

			// Verify the mode was set in the mock client
			updatedMsg, err := mockClient.Connect()
			require.NoError(t, err)

			// Find the heat area and verify mode
			var found bool
			for _, heatArea := range updatedMsg.Device.HeatAreas {
				if heatArea.Nr == roomNr {
					assert.Equal(t, tc.expectedMode, heatArea.Mode, "Mode should be updated to %s (%d)", tc.mode, tc.expectedMode)
					found = true
					break
				}
			}
			assert.True(t, found, "Heat area %d should exist", roomNr)
		})
	}

	// Verify initial mode is still accessible (restore if needed)
	t.Logf("Initial mode was: %d", initialMode)
}

func TestE2E_FullIntegration(t *testing.T) {
	// This test combines both scenarios:
	// 1. Temperature data is published to MQTT
	// 2. Target temperature can be set via MQTT
	// 3. The new target temperature is then published back to MQTT

	// Start MQTT broker
	broker, brokerURL := mqtt.NewBroker(t)
	go func() {
		err := broker.Serve()
		if err != nil {
			t.Logf("broker serve error: %v", err)
		}
	}()
	defer func() {
		err := broker.Close()
		if err != nil {
			t.Logf("broker close error: %v", err)
		}
	}()

	// Wait for broker to be ready
	time.Sleep(100 * time.Millisecond)

	const (
		deviceName = "test-device"
		mqttPrefix = "ezr"
		pollEvery  = 400 * time.Millisecond
		roomNr     = 1
	)

	// Create mock EZR client
	mockClient := mock.NewMockClient()

	// Create store
	memStore := store.NewInMemoryStore()

	// Get initial state
	initialMsg, err := mockClient.Connect()
	require.NoError(t, err)
	deviceID := initialMsg.Device.ID
	memStore.SetID(deviceName, deviceID)

	// Get initial target temperature for room 1
	var initialTargetTemp float64
	for _, ha := range initialMsg.Device.HeatAreas {
		if ha.Nr == roomNr {
			initialTargetTemp = ha.TTarget
			break
		}
	}

	// Create MQTT emitter
	emitter := mqtt.NewEmitter(
		mqtt.WithMqttBrokerUrl[mqtt.Emitter](brokerURL),
		mqtt.WithMqttPrefix[mqtt.Emitter](mqttPrefix),
	)

	// Create handler router
	clients := map[string]transport.Client{
		deviceName: mockClient,
	}
	handlerRouter := handlers.NewHandlerRouter(clients, memStore)

	// Create MQTT listener
	listener := mqtt.NewListener(
		mqtt.WithMqttBrokerUrl[mqtt.Listener](brokerURL),
		mqtt.WithMqttPrefix[mqtt.Listener](mqttPrefix),
	)

	// Connect listener
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	listenerConn, err := listener.Connect(ctx, handlerRouter)
	require.NoError(t, err)
	defer func() {
		err := listenerConn.Disconnect(context.Background())
		if err != nil {
			t.Logf("listener disconnect error: %v", err)
		}
	}()

	// Create test MQTT client for publishing and subscribing
	messageChan := make(chan *paho.Publish, 10)
	router := paho.NewStandardRouter()

	stateTopic := fmt.Sprintf("%s/%s/+/state/+", mqttPrefix, deviceName)
	router.RegisterHandler(stateTopic, func(p *paho.Publish) {
		messageChan <- p
	})

	testClient, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{brokerURL},
		KeepAlive:         10,
		ConnectRetryDelay: 1 * time.Second,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			_, err := manager.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{{Topic: stateTopic}},
			})
			if err != nil {
				t.Logf("subscription error: %v", err)
			}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: "test-integration-client",
			Router:   router,
		},
	})
	require.NoError(t, err)
	defer func() {
		err := testClient.Disconnect(context.Background())
		if err != nil {
			t.Logf("test client disconnect error: %v", err)
		}
	}()

	err = testClient.AwaitConnection(ctx)
	require.NoError(t, err)

	// Give everything time to connect
	time.Sleep(300 * time.Millisecond)

	// Start poller
	poller := polling.NewPoller(deviceName, mockClient, emitter, pollEvery, memStore)
	pollerCtx := t.Context()

	poller.Run(pollerCtx)

	// Wait for initial poll messages
	time.Sleep(600 * time.Millisecond)

	// Clear message channel
	for len(messageChan) > 0 {
		<-messageChan
	}

	// Step 1: Set a new target temperature via MQTT
	newTargetTemp := initialTargetTemp + 2.5

	setTopic := fmt.Sprintf("%s/%s/%d/set/temperature_target", mqttPrefix, deviceName, roomNr)
	_, err = testClient.Publish(ctx, &paho.Publish{
		Topic:   setTopic,
		Payload: []byte(api.FormatFloat(newTargetTemp)),
	})
	require.NoError(t, err)

	// Step 2: Wait for the next poll cycle to publish the updated temperature
	// The poller should pick up the new temperature and publish it
	timeout := time.After(2 * time.Second)
	var foundUpdatedTemp bool

	expectedTopic := fmt.Sprintf("%s/%s/%d/state/temperature_target", mqttPrefix, deviceName, roomNr)

	for !foundUpdatedTemp {
		select {
		case msg := <-messageChan:
			if msg.Topic == expectedTopic {
				if string(msg.Payload) == api.FormatFloat(newTargetTemp) {
					foundUpdatedTemp = true
					t.Logf("Successfully received updated temperature: %.1f", newTargetTemp)
				}
			}
		case <-timeout:
			t.Fatalf("Timeout waiting for updated temperature to be published")
		}
	}

	assert.True(t, foundUpdatedTemp, "Should receive updated target temperature via MQTT")

	// Step 3: Verify the mock client has the updated temperature
	finalMsg, err := mockClient.Connect()
	require.NoError(t, err)

	for _, ha := range finalMsg.Device.HeatAreas {
		if ha.Nr == roomNr {
			assert.Equal(t, newTargetTemp, ha.TTarget,
				"Mock client should have the updated target temperature")
			break
		}
	}
}
