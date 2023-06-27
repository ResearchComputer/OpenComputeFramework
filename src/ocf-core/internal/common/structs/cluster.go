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

type NodeTable struct {
	Nodes []NodeStatus `json:"nodes"`
}

func (lnt NodeTable) Update(node NodeStatus) *NodeTable {
	for idx, n := range lnt.Nodes {
		if n.ClientID == node.ClientID && n.PeerID == node.PeerID {
			if node.Status == "disconnected" {
				lnt.Nodes = append(lnt.Nodes[:idx], lnt.Nodes[idx+1:]...)
				return &lnt
			} else if node.Status == "connected" {
				lnt.Nodes[idx] = node
			}
			return &lnt
		}
	}
	if node.Status == "connected" {
		lnt.Nodes = append(lnt.Nodes, node)
	}
	return &lnt
}

func (lnt NodeTable) FindProviders(service string) []string {
	var providers []string
	for _, n := range lnt.Nodes {
		if n.Service == service {
			providers = append(providers, n.PeerID)
		}
	}
	return providers
}
