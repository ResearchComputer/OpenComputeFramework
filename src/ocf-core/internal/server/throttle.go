package server

import (
	"ocfcore/internal/cluster"
	"ocfcore/internal/common/structs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var instructionsHub structs.WorkloadInstructionsHub

func LoadWorkload(g *gin.Context) {
	workerID := g.Param("workerId")
	var provisionInstructions structs.ProvisionModelsPlan
	err := g.BindJSON(&provisionInstructions)
	if err != nil {
		g.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if instructionsHub.Instructions == nil {
		instructionsHub.Instructions = make(map[string]structs.ProvisionModelsPlan)
	}
	// check if the worker is already connected
	for _, worker := range workerHub.Workers {
		if worker.WorkerID == workerID {
			// update the worker
			instructionsHub.Instructions[workerID] = provisionInstructions
			g.JSON(200, gin.H{"message": "ok"})
			return
		}
	}
	g.JSON(404, gin.H{"error": "worker not found"})
}

func GetWorkloadInstructions(g *gin.Context) {
	workerID := g.Param("workerId")
	g.JSON(200, instructionsHub.Instructions[workerID])
}

func AddClusterNode(g *gin.Context) {
	var acquirePayload map[string]string
	g.BindJSON(&acquirePayload)
	acquireMachinePayload := structs.AcquireMachinePayload{
		Script: viper.GetString("acquire_machine.script"),
		Params: acquirePayload,
	}
	cluster.NewSlurmClusterClient().AcquireMachine(acquireMachinePayload)
}
