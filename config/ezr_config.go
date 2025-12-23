package config

type HttpClientConfig struct {
	Host string `mapstructure:"host" json:"host" validate:"required"`
}

type EzrConfig struct {
	Name string            `mapstructure:"name" json:"name" validate:"required"`
	Type string            `mapstructure:"type" toml:"type" validate:"required,oneof=http mock"`
	Http *HttpClientConfig `mapstructure:"http,omitempty" toml:"http,omitempty" validate:"required_if=Type http"`
}
