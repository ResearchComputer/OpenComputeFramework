package structs

import (
	"ocfcore/internal/protocol/p2p"
	"time"
)

// Cluster is a lower-level interface from workers

type AcquireMachinePayload struct {
	Script string            `json:"script"`
	Params map[string]string `json:"params"`
}

type WarmedMachine struct {
	MachineID string        `json:"machine_id"`
	Status    string        `json:"status"`
	Life      time.Duration `json:"life"`
	StartedAt int64         `json:"started_at"`
}

type NodeStatus struct {
	PeerID   string        `json:"peer_id"`
	ClientID int           `json:"client_id"`
	Status   string        `json:"status"`
	Specs    []p2p.GPUSpec `json:"gpus"`
	Service  string        `json:"service"`
}
