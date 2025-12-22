package config

type GeneralConfig struct {
	PollEvery string `mapstructure:"poll_every" json:"poll_every" validate:"required,gt=0"`
}
