package cmd

import (
	"ocfcore/internal/daemon"

	"github.com/spf13/cobra"
)

var starocfcored = &cobra.Command{
	Use:   "start",
	Short: "Start pulling jobs from the zentrum and dispatching them to the workers",
	Run: func(cmd *cobra.Command, args []string) {
		daemon.Start()
	}}
