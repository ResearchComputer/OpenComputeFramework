package protocol

import (
	"context"
	"encoding/json"
	"ocf/internal/common"
	"sync"

	ds "github.com/ipfs/go-datastore"
)

var dntOnce sync.Once

const (
	CONNECTED    string = "connected"
	DISCONNECTED string = "disconnected"
)

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

type Service struct {
	Name     string         `json:"name"`
	Hardware []HardwareSpec `json:"hardware"`
	Status   string         `json:"status"`
}

// Peer is a single node in the network, as can be seen by the current node.
type Peer struct {
	Latency           int       `json:"latency"` // in ms
	Privileged        bool      `json:"privileged"`
	Owner             string    `json:"owner"`
	CurrentOffering   []string  `json:"current_offering"`
	Role              []string  `json:"role"`
	Status            string    `json:"status"`
	AvailableOffering []string  `json:"available_offering"`
	Service           []Service `json:"service"`
	LastSeen          int64     `json:"last_seen"`
	Version           string    `json:"version"`
}

// Node table tracks the nodes and their status in the network.
// This is also a
type NodeTable map[string]Peer

var dnt *NodeTable

func GetNodeTable() *NodeTable {
	dntOnce.Do(func() {
		dnt = &NodeTable{}
	})
	return dnt
}

func UpdateNodeTable(peer Peer) {
	ctx := context.Background()
	host, _ := GetP2PNode(nil)
	// broadcast the peer to the network
	store, pcancel := GetCRDTStore()
	key := ds.NewKey(host.ID().String())
	value, err := json.Marshal(peer)
	common.ReportError(err, "Error while marshalling peer")
	store.Put(ctx, key, value)
	defer pcancel()
}

func UpdateNodeTableHook(key ds.Key, value []byte) {
	table := *GetNodeTable()
	var peer Peer
	err := json.Unmarshal(value, &peer)
	common.ReportError(err, "Error while unmarshalling peer")
	table[key.String()] = peer
}
