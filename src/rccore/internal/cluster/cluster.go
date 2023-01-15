package cluster

import "rccore/internal/common/structs"

type ClusterManager interface {
	AcquireMachine(payload structs.AcquireMachinePayload)
	Execute(command string)
}
