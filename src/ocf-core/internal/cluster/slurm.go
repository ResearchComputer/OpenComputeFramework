package cluster

import (
	"os/exec"
	"ocfcore/internal/common"
	"ocfcore/internal/common/structs"
	"strings"
	"time"
)

type SlurmCluster struct {
	ConnectedMachine []structs.WarmedMachine
}

var slurmClusterClient *SlurmCluster

func NewSlurmClusterClient() *SlurmCluster {
	if slurmClusterClient == nil {
		slurmClusterClient = &SlurmCluster{}
	}
	return slurmClusterClient
}

func (s *SlurmCluster) AcquireMachine(payload structs.AcquireMachinePayload) {
	common.Logger.Info("Acquiring machine", "payload", payload)
	output, err := s.execute(payload.Script)
	if err != nil {
		common.Logger.Error("Could not acquire machine", "error", err)
		return
	}
	outputString := strings.Split(output, " ")
	machineID := outputString[len(outputString)-1]
	common.Logger.Info("Machine ID", "machineID", machineID)
	s.ConnectedMachine = append(s.ConnectedMachine, structs.WarmedMachine{
		MachineID: machineID,
		Status:    "REQUESTING",
		StartedAt: time.Now().Unix(),
		Life:      time.Hour * 4,
	})
}

func (s *SlurmCluster) execute(command string) (string, error) {
	// execute the command
	prg := "sbatch"
	args := []string{command}
	cmd := exec.Command(prg, args...)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}
