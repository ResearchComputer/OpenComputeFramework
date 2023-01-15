package cmd

import (
	"rccore/internal/daemon"

	"github.com/spf13/cobra"
)

var starrccored = &cobra.Command{
	Use:   "start",
	Short: "Start pulling jobs from the zentrum and dispatching them to the workers",
	Run: func(cmd *cobra.Command, args []string) {
		daemon.Start()
	}}
