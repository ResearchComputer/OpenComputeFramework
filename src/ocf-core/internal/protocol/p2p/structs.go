package p2p

import (
	"ocfcore/internal/common"
	"sync"
	"time"
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
	PeerID            string    `json:"peer_id"`
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
type NodeTable struct {
	Peers []Peer `json:"peers"`
}

var nodeTable *NodeTable

func GetNodeTable() *NodeTable {
	dntOnce.Do(func() {
		nodeTable = &NodeTable{Peers: []Peer{}}
	})
	return nodeTable
}

func (dnt *NodeTable) Update(peer Peer) *NodeTable {
	for idx, n := range dnt.Peers {
		if n.PeerID == peer.PeerID {
			dnt.Peers[idx].LastSeen = time.Now().Unix()
			if peer.Status == DISCONNECTED {
				dnt.Peers[idx].Status = DISCONNECTED
				dnt.Peers[idx].CurrentOffering = []string{}
				dnt.Peers[idx].Service = []Service{}
				dnt.Peers[idx].LastSeen = time.Now().Unix()
				return dnt
			}
			return dnt
		}
	}
	if peer.Status == CONNECTED {
		dnt.Peers = append(dnt.Peers, peer)
	}
	return dnt
}

func (dnt NodeTable) FindProviders(service string) []Peer {
	var providers []Peer
	for _, p := range dnt.Peers {
		for _, s := range p.Service {
			if s.Name == service {
				providers = append(providers, p)
			}
		}
	}
	return providers
}

func (dnt *NodeTable) RemoveDisconnectedPeers(disconnected []string) {
	for _, p := range dnt.Peers {
		for _, d := range disconnected {
			if p.PeerID == d {
				dnt.Update(Peer{PeerID: p.PeerID, Status: DISCONNECTED})
			}
		}
	}
}

func (dnt *NodeTable) NewOffering(peerId string, newService string, hardware HardwareSpec) {
	for idx, p := range dnt.Peers {
		if p.PeerID == peerId {
			dnt.Peers[idx].CurrentOffering = append(dnt.Peers[idx].CurrentOffering, newService)
			dnt.Peers[idx].Service = append(dnt.Peers[idx].Service, Service{Name: newService, Hardware: []HardwareSpec{hardware}})
			BroadcastPeerOffering(dnt.Peers[idx])
			break
		}
	}
}

func (dnt *NodeTable) RemoveOffering(peerId string, newService string) {
	for idx, p := range dnt.Peers {
		if p.PeerID == peerId {
			dnt.Peers[idx].CurrentOffering = common.RemoveString(dnt.Peers[idx].CurrentOffering, newService)
			// remove service from service list
			for i, s := range dnt.Peers[idx].Service {
				if s.Name == newService {
					dnt.Peers[idx].Service = append(dnt.Peers[idx].Service[:i], dnt.Peers[idx].Service[i+1:]...)
					break
				}
			}
			BroadcastPeerOffering(dnt.Peers[idx])
			break
		}
	}
}

func (dnt *NodeTable) UpdateNodeTable(peer Peer) {
	for idx, p := range dnt.Peers {
		if p.PeerID == peer.PeerID {
			dnt.Peers[idx] = peer
			break
		}
	}
}
