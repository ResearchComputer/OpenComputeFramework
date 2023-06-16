package cmd

import (
	"ocfcore/internal/cluster"
	"ocfcore/internal/common"
	"ocfcore/internal/common/structs"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var acquireCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a machine to the cluster effectively",
	Run: func(cmd *cobra.Command, args []string) {
		slurmClient := cluster.NewSlurmClusterClient()
		slurmClient.AcquireMachine(structs.AcquireMachinePayload{
			Script: viper.GetString("acquire_machine.script"),
			Params: make(map[string]string, 0),
		})
	},
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage the cluster status",
	Long:  `Manage the cluster status, currently supporting slurm, kubernetes and baremetal`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			common.Logger.Error("Could not print help", "error", err)
		}
	},
}
