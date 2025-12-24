package mqtt

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type Listener struct {
	connectionDetails
	mqttGroup string
}

func NewListener(opts ...Opt[Listener]) *Listener {
	l := new(Listener)
	for _, opt := range opts {
		opt(l)
	}
	ensureListenerDefaults(l)
	return l
}

func ensureListenerDefaults(l *Listener) {
	if l.mqttBrokerUrls == nil {
		u, err := url.Parse("mqtt://127.0.0.1:1883/")
		if err != nil {
			panic(err)
		}
		l.mqttBrokerUrls = []*url.URL{u}
	}
	if l.mqttPrefix == "" {
		l.mqttPrefix = "ezr"
	}
	if l.mqttGroup == "" {
		l.mqttGroup = "ezr2mqtt"
	}
	if l.mqttConnectTimeout == 0 {
		l.mqttConnectTimeout = 10 * time.Second
	}
	if l.mqttConnectRetryDelay == 0 {
		l.mqttConnectRetryDelay = 1 * time.Second
	}
	if l.mqttKeepAliveInterval == 0 {
		l.mqttKeepAliveInterval = 10
	}
}

func (l *Listener) Connect(ctx context.Context, handler api.MessageHandler) (api.Connection, error) {
	var err error

	ctx, cancel := context.WithTimeout(ctx, l.mqttConnectTimeout)
	defer cancel()

	clientId := fmt.Sprintf("%s-%s", l.mqttGroup, randSeq(5))

	readyCh := make(chan struct{})

	// e.g. ezr/name123/bedroom/set/temperature
	topic := fmt.Sprintf("%s/+/+/set/+", l.mqttPrefix)

	conn := new(connection)
	mqttRouter := paho.NewStandardRouter()
	conn.mqttConn, err = autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		ServerUrls:        l.mqttBrokerUrls,
		ConnectUsername:   l.mqttUsername,
		ConnectPassword:   []byte(l.mqttPassword),
		KeepAlive:         l.mqttKeepAliveInterval,
		ConnectRetryDelay: l.mqttConnectRetryDelay,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			_, err := manager.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{{Topic: topic}},
			})
			if err != nil {
				slog.Error("failed to subscribe to topic", "topic", topic)
				return
			}
			mqttRouter.UnregisterHandler(topic)
			mqttRouter.RegisterHandler(topic, func(mqttMsg *paho.Publish) {
				ctx := context.Background()

				// determine parts - ezr/name123/bedroom/set/temperature
				topicParts := strings.Split(mqttMsg.Topic, "/")

				name := topicParts[len(topicParts)-4]
				t := topicParts[len(topicParts)-1]
				room, err := strconv.Atoi(topicParts[len(topicParts)-3])
				if err != nil {
					slog.Error("invalid room number in topic", "topic", mqttMsg.Topic)
				}

				msg := api.Message{
					Room: room,
					Type: t,
					Data: string(mqttMsg.Payload),
				}

				// execute the handler
				handler.Handle(ctx, name, &msg)
			})
			readyCh <- struct{}{}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: clientId,
			Router:   mqttRouter,
		},
	})
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, errors.New("timeout waiting for mqtt connectionDetails setup")
	case <-readyCh:
		return conn, nil
	}
}

type connection struct {
	mqttConn *autopaho.ConnectionManager
}

func (c *connection) Disconnect(ctx context.Context) error {
	if c.mqttConn != nil {
		err := c.mqttConn.Disconnect(ctx)
		if err != nil {
			return err
		}
		c.mqttConn = nil
	}
	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		//#nosec G404 - client suffix does not require secure random number generator
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
