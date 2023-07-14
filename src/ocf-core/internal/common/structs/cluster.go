package structs

import "time"

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

type GPUSpec struct {
	Name       string `json:"name"`
	Memory     int64  `json:"memory"`
	FreeMemory int64  `json:"memory_free"`
	UsedMemory int64  `json:"memory_used"`
}

type NodeStatus struct {
	PeerID   string    `json:"peer_id"`
	ClientID int       `json:"client_id"`
	Status   string    `json:"status"`
	Specs    []GPUSpec `json:"gpus"`
	Service  string    `json:"service"`
}
