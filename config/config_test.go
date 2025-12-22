package config_test

import (
	"testing"

	"github.com/chrishrb/ezr2mqtt/config"
	clone "github.com/huandu/go-clone/generic"
	"github.com/stretchr/testify/require"
)

func TestConfigure(t *testing.T) {
	cfg := clone.Clone(&config.DefaultConfig)

	_, err := config.Configure(t.Context(), cfg)
	require.NoError(t, err)
}
