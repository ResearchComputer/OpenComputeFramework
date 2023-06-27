package cmd

import (
	"fmt"
	"net/http"

	"ocfcore/internal/common"

	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the ocf binary to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		updateURL := "https://cdn.xzyao.dev/ocfcore"
		resp, err := http.Get(updateURL)
		if err != nil {
			common.Logger.Error("Error while checking for updates: ", err)
		}
		defer resp.Body.Close()
		err = selfupdate.Apply(resp.Body, selfupdate.Options{})
		if err != nil {
			if rerr := selfupdate.RollbackError(err); rerr != nil {
				common.Logger.Info("Failed to rollback from bad update: ", rerr)
			}
		}
		common.Logger.Info("Successfully updated")
		fmt.Printf("ocfcore version %s", common.JSONVersion.Version)
		fmt.Printf(" (commit: %s)", common.JSONVersion.Commit)
		fmt.Printf(" (built at: %s)", common.JSONVersion.Date)
		fmt.Println()
	},
}
