package config

type EzrConfig struct {
	Host string `mapstructure:"host" json:"host" validate:"required"`
}
