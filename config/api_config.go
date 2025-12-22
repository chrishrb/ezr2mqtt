package config

type MqttSettingsConfig struct {
	Urls              []string `mapstructure:"urls" toml:"urls" validate:"required,dive,required"`
	Prefix            string   `mapstructure:"prefix" toml:"prefix" validate:"required"`
	Group             string   `mapstructure:"group" toml:"group" validate:"required"`
	ConnectTimeout    string   `mapstructure:"connect_timeout" toml:"connect_timeout" validate:"required"`
	ConnectRetryDelay string   `mapstructure:"connect_retry_delay" toml:"connect_retry_delay" validate:"required"`
	KeepAliveInterval string   `mapstructure:"keep_alive_interval" toml:"keep_alive_interval" validate:"required"`
}

type ApiSettingsConfig struct {
	Mqtt *MqttSettingsConfig `mapstructure:"mqtt,omitempty" json:"mqtt,omitempty"`
}
