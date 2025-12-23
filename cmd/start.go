package cmd

import (
	"context"
	"log/slog"

	"github.com/chrishrb/ezr2mqtt/config"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Run the ezr2mqtt service",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.DefaultConfig
		if configFile != "" {
			err := cfg.LoadFromFile(configFile)
			if err != nil {
				return err
			}
		}

		settings, err := config.Configure(context.Background(), &cfg)
		if err != nil {
			return err
		}

		errCh := make(chan error, 1)

		// Connect to mqtt broker and start listening for messages
		conn, err := settings.MqttListener.Connect(context.Background(), settings.MqttHandler)
		if err != nil {
			errCh <- err
		}

		// Start periodic requests
		periodicRequester := settings.PeriodicRequester
		for _, pr := range periodicRequester {
			pr.Run(context.Background())
		}

		slog.Info("ezr2mqtt started")

		err = <-errCh

		if conn != nil {
			err := conn.Disconnect(context.Background())
			if err != nil {
				slog.Warn("closing transport connection", "error", err)
			}
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&configFile, "config-file", "c", "/config/config.yaml",
		"The config file to use")
}
