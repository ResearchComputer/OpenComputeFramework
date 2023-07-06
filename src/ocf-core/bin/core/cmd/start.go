package cmd

import (
	"ocfcore/internal/daemon"

	"github.com/spf13/cobra"
)

var starocfcored = &cobra.Command{
	Use:   "start",
	Short: "Start listening for incoming connections",
	Run: func(cmd *cobra.Command, args []string) {
		daemon.Start()
	}}
