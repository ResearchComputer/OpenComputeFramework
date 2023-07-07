package p2p

type GPUSpec struct {
	Name            string `json:"name"`
	Memory          int64  `json:"memory"`
	MemoryBandwidth int64  `json:"memory_bandwidth"`
	UsedMemory      int64  `json:"memory_used"`
}

type HardwareSpec struct {
	GPUs            []GPUSpec `json:"gpus"`
	Memory          int64     `json:"host_memory"`
	MemoryBandwidth int64     `json:"host_memory_bandwidth"`
	UsedMemory      int64     `json:"host_memory_used"`
}

type PeerStatus string

const (
	CONNECTED    PeerStatus = "connected"
	DISCONNECTED PeerStatus = "disconnected"
)

// Peer is a single node in the network, as can be seen by the current node.
type Peer struct {
	PeerID          string         `json:"peer_id"`
	Latency         int            `json:"latency"` // in ms
	Privileged      bool           `json:"privileged"`
	Owner           string         `json:"owner"`
	CurrentOffering []string       `json:"current_offering"`
	Hardware        []HardwareSpec `json:"hardware"`
	Role            []string       `json:"role"`
	Status          string         `json:"status"`
	Service         string         `json:"service"`
}

// Node table tracks the nodes and their status in the network.
type NodeTable struct {
	Peers []Peer `json:"peers"`
}

func (dnt NodeTable) Update(peer Peer) *NodeTable {
	for idx, n := range dnt.Peers {
		if n.PeerID == peer.PeerID {
			if peer.Status == "disconnected" {
				dnt.Peers = append(dnt.Peers[:idx], dnt.Peers[idx+1:]...)
				return &dnt
			} else if peer.Status == "connected" {
				dnt.Peers[idx] = peer
			}
			return &dnt
		}
	}
	if peer.Status == "connected" {
		dnt.Peers = append(dnt.Peers, peer)
	}
	return &dnt
}

func (dnt NodeTable) FindProviders(service string) []Peer {
	var providers []Peer
	for _, n := range dnt.Peers {
		if n.Service == service {
			providers = append(providers, n)
		}
	}
	return providers
}
