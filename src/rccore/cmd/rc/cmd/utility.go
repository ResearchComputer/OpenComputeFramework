package cmd

import (
	"fmt"
	"rccore/internal/common"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of rccore",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("rccore version %s", common.JSONVersion.Version)
		fmt.Printf(" (commit: %s)", common.JSONVersion.Commit)
		fmt.Printf(" (built at: %s)", common.JSONVersion.Date)
		fmt.Println()
	},
}
