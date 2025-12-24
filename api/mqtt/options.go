package mqtt

import (
	"net/url"
	"time"
)

type connectionDetails struct {
	mqttBrokerUrls        []*url.URL
	mqttUsername          string
	mqttPassword          string
	mqttPrefix            string
	mqttConnectTimeout    time.Duration
	mqttConnectRetryDelay time.Duration
	mqttKeepAliveInterval uint16
}

type Opt[T any] func(h *T)

func WithMqttBrokerUrl[T Emitter | Listener](brokerUrl *url.URL) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Emitter:
			x.mqttBrokerUrls = append(x.mqttBrokerUrls, brokerUrl)
		case *Listener:
			x.mqttBrokerUrls = append(x.mqttBrokerUrls, brokerUrl)
		}
	}
}

func WithMqttBrokerUrls[T Emitter | Listener](brokerUrls []*url.URL) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Emitter:
			x.mqttBrokerUrls = brokerUrls
		case *Listener:
			x.mqttBrokerUrls = brokerUrls
		}
	}
}

func WithMqttPrefix[T Emitter | Listener](mqttPrefix string) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Emitter:
			x.mqttPrefix = mqttPrefix
		case *Listener:
			x.mqttPrefix = mqttPrefix
		}
	}
}

func WithMqttConnectSettings[T Emitter | Listener](mqttConnectTimeout, mqttConnectRetryDelay, mqttKeepAliveInterval time.Duration) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Emitter:
			x.mqttConnectTimeout = mqttConnectTimeout
			x.mqttConnectRetryDelay = mqttConnectRetryDelay
			x.mqttKeepAliveInterval = uint16(mqttKeepAliveInterval.Round(time.Second).Seconds())
		case *Listener:
			x.mqttConnectTimeout = mqttConnectTimeout
			x.mqttConnectRetryDelay = mqttConnectRetryDelay
			x.mqttKeepAliveInterval = uint16(mqttKeepAliveInterval.Round(time.Second).Seconds())
		}
	}
}

func WithMqttUsername[T Emitter | Listener](mqttUsername string) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Emitter:
			x.mqttUsername = mqttUsername
		case *Listener:
			x.mqttUsername = mqttUsername
		}
	}
}

func WithMqttPassword[T Emitter | Listener](mqttPassword string) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Emitter:
			x.mqttPassword = mqttPassword
		case *Listener:
			x.mqttPassword = mqttPassword
		}
	}
}

func WithMqttGroup[T Listener](mqttGroup string) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Listener:
			x.mqttGroup = mqttGroup
		}
	}
}
