package mqtt_test

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	mqtt2 "github.com/chrishrb/ezr2mqtt/api/mqtt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmitterSendsMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// start the broker
	broker, clientUrl := mqtt2.NewBroker(t)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()

	err := broker.Serve()
	require.NoError(t, err)

	emitter := mqtt2.NewEmitter(
		mqtt2.WithMqttBrokerUrl[mqtt2.Emitter](clientUrl),
		mqtt2.WithMqttPrefix[mqtt2.Emitter]("ezr"))

	// subscribe to the output channel
	rcvdCh := make(chan struct{})

	router := paho.NewStandardRouter()
	router.RegisterHandler("ezr/id123/1/state/temperature", func(publish *paho.Publish) {
		assert.Equal(t, "ezr/id123/1/state/temperature", publish.Topic)
		var msg api.Message
		err := json.Unmarshal(publish.Payload, &msg)
		assert.NoError(t, err, "payload is not the expected message type")
		assert.Equal(t, 332, 332)
		rcvdCh <- struct{}{}
	})
	mqttClient := listenForMessageSentByManager(t, ctx, clientUrl, router)

	defer func() {
		_ = mqttClient.Disconnect(ctx)
	}()

	// publish a message to the input channel
	msg := api.Message{
		Room: 1,
		Type: "temperature",
		Data: 22,
	}
	err = emitter.Emit(context.Background(), "id123", &msg)
	require.NoError(t, err)

	// wait for success
	select {
	case <-rcvdCh:
		// success
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test")
	}
}

func listenForMessageSentByManager(t *testing.T, ctx context.Context, clientUrl *url.URL, router paho.Router) *autopaho.ConnectionManager {
	mqttClient, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{clientUrl},
		KeepAlive:         10,
		ConnectRetryDelay: 10,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			_, err := manager.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{
					{
						Topic: "ezr/#",
					},
				},
			})
			require.NoError(t, err)
		},
		ClientConfig: paho.ClientConfig{
			ClientID: "test",
			Router:   router,
		},
	})
	require.NoError(t, err)

	err = mqttClient.AwaitConnection(ctx)
	require.NoError(t, err)

	return mqttClient
}
