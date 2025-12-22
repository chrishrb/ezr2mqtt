package config

import (
	"bufio"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert/yaml"
)

type BaseConfig struct {
	Api ApiSettingsConfig `mapstructure:"api" json:"api" validate:"required"`
	// TODO: make possible to have multiple transport configs
	Ezr     EzrConfig     `mapstructure:"ezr" json:"ezr" validate:"required"`
	General GeneralConfig `mapstructure:"general" json:"general" validate:"required"`
}

// DefaultConfig provides the default configuration. The configuration
// read from the YAML file will overlay this configuration.
var DefaultConfig = BaseConfig{
	General: GeneralConfig{
		PollEvery: "1m",
	},
}

// Load reads YAML configuration from a reader.
func (c *BaseConfig) Load(reader io.Reader) error {
	b, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, c); err != nil {
		return err
	}
	return nil
}

// LoadFromFile reads YAML configuration from a file.
func (c *BaseConfig) LoadFromFile(configFile string) error {
	//#nosec G304 - only files specified by the person running the application will be loaded
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}
	err = c.Load(bufio.NewReader(f))
	return err
}

// Validate ensures that the configuration is structurally valid.
func (c *BaseConfig) Validate() error {
	validate := validator.New()

	return validate.Struct(c)
}
