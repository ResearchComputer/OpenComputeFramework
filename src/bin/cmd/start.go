package cmd

import (
	"ocf/internal/server"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start listening for incoming connections",
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer()
	}}
