package config

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/api/mqtt"
	"github.com/chrishrb/ezr2mqtt/handlers"
	"github.com/chrishrb/ezr2mqtt/periodic"
	"github.com/chrishrb/ezr2mqtt/transport"
	"github.com/chrishrb/ezr2mqtt/transport/http"
	"github.com/chrishrb/ezr2mqtt/transport/mock"
)

type Config struct {
	EzrClient         transport.Client
	MqttListener      api.Listener
	MqttEmitter       api.Emitter
	MqttHandler       api.MessageHandler
	PeriodicRequester *periodic.PeriodicRequester
}

func Configure(ctx context.Context, cfg *BaseConfig) (c *Config, err error) {
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	c = &Config{}

	c.EzrClient, err = getEzrClient(cfg.Ezr)
	if err != nil {
		return nil, err
	}

	c.MqttListener, err = getMqttReceiver(cfg.Api)
	if err != nil {
		return nil, err
	}

	c.MqttEmitter, err = getMqttEmitter(cfg.Api)
	if err != nil {
		return nil, err
	}

	c.MqttHandler = handlers.NewHandlerRouter(c.EzrClient)

	c.PeriodicRequester, err = getPeriodicRequester(c.EzrClient, c.MqttEmitter, cfg.General)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getEzrClient(cfg EzrConfig) (transport.Client, error) {
	switch cfg.Type {
	case "http":
		return http.NewHTTPClient(cfg.Http.Host), nil
	case "mock":
		return mock.NewMockClient(), nil
	default:
		return nil, fmt.Errorf("unsupported ezr client type: %s", cfg.Type)
	}
}

func getMqttReceiver(cfg ApiSettingsConfig) (api.Listener, error) {
	switch cfg.Type {
	case "mqtt":
		var mqttUrls []*url.URL
		for _, urlStr := range cfg.Mqtt.Urls {
			u, err := url.Parse(urlStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mqtt url: %w", err)
			}
			mqttUrls = append(mqttUrls, u)
		}

		mqttConnectTimeout, err := time.ParseDuration(cfg.Mqtt.ConnectTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect timeout: %w", err)
		}

		mqttConnectRetryDelay, err := time.ParseDuration(cfg.Mqtt.ConnectRetryDelay)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect retry delay: %w", err)
		}

		mqttKeepAliveInterval, err := time.ParseDuration(cfg.Mqtt.KeepAliveInterval)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt keep alive interval: %w", err)
		}

		opts := []mqtt.Opt[mqtt.Listener]{
			mqtt.WithMqttBrokerUrls[mqtt.Listener](mqttUrls),
			mqtt.WithMqttPrefix[mqtt.Listener](cfg.Mqtt.Prefix),
			mqtt.WithMqttConnectSettings[mqtt.Listener](mqttConnectTimeout, mqttConnectRetryDelay, mqttKeepAliveInterval),
			mqtt.WithMqttGroup(cfg.Mqtt.Group),
		}

		return mqtt.NewListener(opts...), nil
	default:
		return nil, fmt.Errorf("unsupported api type: %s", cfg.Type)
	}
}
func getMqttEmitter(cfg ApiSettingsConfig) (api.Emitter, error) {
	switch cfg.Type {
	case "mqtt":
		var mqttUrls []*url.URL
		for _, urlStr := range cfg.Mqtt.Urls {
			u, err := url.Parse(urlStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mqtt url: %w", err)
			}
			mqttUrls = append(mqttUrls, u)
		}

		mqttConnectTimeout, err := time.ParseDuration(cfg.Mqtt.ConnectTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect timeout: %w", err)
		}

		mqttConnectRetryDelay, err := time.ParseDuration(cfg.Mqtt.ConnectRetryDelay)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect retry delay: %w", err)
		}

		mqttKeepAliveInterval, err := time.ParseDuration(cfg.Mqtt.KeepAliveInterval)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt keep alive interval: %w", err)
		}

		mqttEmitter := mqtt.NewEmitter(
			mqtt.WithMqttBrokerUrls[mqtt.Emitter](mqttUrls),
			mqtt.WithMqttPrefix[mqtt.Emitter](cfg.Mqtt.Prefix),
			mqtt.WithMqttConnectSettings[mqtt.Emitter](mqttConnectTimeout, mqttConnectRetryDelay, mqttKeepAliveInterval))

		return mqttEmitter, nil
	default:
		return nil, fmt.Errorf("unsupported api type: %s", cfg.Type)
	}
}

func getPeriodicRequester(client transport.Client, emitter api.Emitter, cfg GeneralConfig) (*periodic.PeriodicRequester, error) {
	runEvery, err := time.ParseDuration(cfg.PollEvery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse periodic PollEvery: %w", err)
	}

	return periodic.NewPeriodicRequester(client, emitter, runEvery), nil
}
