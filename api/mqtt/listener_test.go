package mqtt_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/api/mqtt"
	server "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListenerProcessesMessagesReceivedFromTheBroker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// start the broker
	broker, clientUrl := mqtt.NewBroker(t)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()
	err := broker.Serve()
	require.NoError(t, err)

	// setup the handler
	receivedMsgCh := make(chan struct{})
	handler := func(ctx context.Context, id string, msg *api.Message) {
		assert.Equal(t, "id123", id)
		assert.Equal(t, 1, msg.Room)
		assert.Equal(t, "temperature", msg.Type)
		assert.Equal(t, 23.2, msg.Data)
		receivedMsgCh <- struct{}{}
	}

	// connect the listener to the broker
	listener := mqtt.NewListener(mqtt.WithMqttBrokerUrl[mqtt.Listener](clientUrl))
	conn, err := listener.Connect(ctx, api.MessageHandlerFunc(handler))
	require.NoError(t, err)
	defer func() {
		if conn != nil {
			err := conn.Disconnect(ctx)
			require.NoError(t, err)
		}
	}()

	// publish message
	publishMessage(t, broker, api.Message{
		Room: 1,
		Type: "temperature",
		Data: 23.2,
	})

	// wait for message to be received / timeout
	select {
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test to complete")
	case <-receivedMsgCh:
		// do nothing
	}
}

func publishMessage(t *testing.T, broker *server.Server, msg api.Message) {
	msgBytes, err := json.Marshal(msg)
	require.NoError(t, err)

	cl := broker.NewClient(nil, "local", "inline", true)
	err = broker.InjectPacket(cl, packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type:   packets.Publish,
			Qos:    0,
			Retain: false,
		},
		TopicName: "ezr/id123/bedroom/set/temperature",
		Payload:   msgBytes,
		PacketID:  uint16(0),
	})
	require.NoError(t, err)
}
