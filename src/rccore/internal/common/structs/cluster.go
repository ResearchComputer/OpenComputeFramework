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
