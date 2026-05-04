/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/diyorbek/sentinel/agent/cmd/app"
	"github.com/spf13/cobra"
)

var (
	apiKey    string
	server    string
	agentName string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start sentinel agent",

	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("starting sentinel agent...")

		// PID yozish
		if err := writePID(); err != nil {
			slog.Warn("could not write pid file", "err", err)
		}
		defer os.Remove(pidFilePath) // agent to'xtaganda o'chiradi

		ctx := context.Background()

		// 1. INIT APP
		a, err := app.New(apiKey, server, agentName)
		if err != nil {
			slog.Error("failed to init app", "err", err)
			return
		}
		defer a.Close()

		// 2. RUN COLLECTORS
		a.RunCollectors(ctx)

		// 3. RUN SENDER
		a.RunSender(ctx)

		a.Heartbeat()

		slog.Info("agent started successfully")

		// 4. BLOCK FOREVER
		select {}
	},
}

func init() {
	startCmd.Flags().StringVar(&apiKey, "api-key", "", "API key for authentication")
	startCmd.Flags().StringVar(&server, "server", "", "Server URL")
	startCmd.Flags().StringVar(&agentName, "name", "", "Agent Server Name")

	startCmd.MarkFlagRequired("api-key")
	startCmd.MarkFlagRequired("server")

	rootCmd.AddCommand(startCmd)
}

// Joriy process PID ni faylga yozadi
func writePID() error {
	pid := strconv.Itoa(os.Getpid())
	return os.WriteFile(pidFilePath, []byte(pid), 0644)
}
