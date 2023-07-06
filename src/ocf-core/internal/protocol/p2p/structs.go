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

// OCFPeer is a single node in the network, as can be seen by the current node.
type OCFPeer struct {
	PeerID          string
	Latency         int // in
	Privileged      bool
	Owner           string
	CurrentOffering []string
	Hardware        []HardwareSpec
	Role            []string
}

// Node table tracks the nodes and their status in the network.
type NodeTable struct {
}
