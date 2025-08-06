package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"ocf/internal/common"
	"ocf/internal/platform"
	"sync"

	ds "github.com/ipfs/go-datastore"
	"github.com/spf13/viper"
)

var dntOnce sync.Once
var myself Peer

const (
	CONNECTED    string = "connected"
	DISCONNECTED string = "disconnected"
)

type Service struct {
	Name     string              `json:"name"`
	Hardware common.HardwareSpec `json:"hardware"`
	Status   string              `json:"status"`
	Host     string              `json:"host"`
	Port     string              `json:"port"`
	// IdentityGroup is a list of identities that can access this service
	// Format: <identity_group_name>=<identity_name>
	// e.g., "model=resnet50"
	IdentityGroup []string `json:"identity_group"`
}

// Peer is a single node in the network, as can be seen by the current node.
type Peer struct {
	ID                string              `json:"id"`
	Latency           int                 `json:"latency"` // in ms
	Privileged        bool                `json:"privileged"`
	Owner             string              `json:"owner"`
	CurrentOffering   []string            `json:"current_offering"`
	Role              []string            `json:"role"`
	Status            string              `json:"status"`
	AvailableOffering []string            `json:"available_offering"`
	Service           []Service           `json:"service"`
	LastSeen          int64               `json:"last_seen"`
	Version           string              `json:"version"`
	PublicAddress     string              `json:"public_address"`
	Hardware          common.HardwareSpec `json:"hardware"`
	Connected         bool                `json:"connected"`
	Load              []int               `json:"load"`
}

type PeerWithStatus struct {
	ID            string `json:"id"`
	Connectedness string `json:"connectedness"` // "connected" or "disconnected"
}

// Node table tracks the nodes and their status in the network.
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
	store, _ := GetCRDTStore()
	key := ds.NewKey(host.ID().String())
	peer.ID = host.ID().String()
	// merge services instead of overwriting
	// first find the peer in the table if it exists
	existingPeer, err := GetPeerFromTable(peer.ID)
	if err == nil {
		peer.Service = append(peer.Service, existingPeer.Service...)
	}
	if viper.GetString("public-addr") != "" {
		peer.PublicAddress = viper.GetString("public-addr")
	}
	value, err := json.Marshal(peer)
	common.ReportError(err, "Error while marshalling peer")
	store.Put(ctx, key, value)
}

func MarkSelfAsBootstrap() {
	if viper.GetString("public-addr") != "" {
		common.Logger.Info("Registering myself as a bootstrap node")
		ctx := context.Background()
		store, _ := GetCRDTStore()
		host, _ := GetP2PNode(nil)
		key := ds.NewKey(host.ID().String())
		peer := Peer{
			ID:            host.ID().String(),
			PublicAddress: viper.GetString("public-addr"),
		}
		value, err := json.Marshal(peer)
		common.ReportError(err, "Error while marshalling peer")
		store.Put(ctx, key, value)
	}
}

func DeleteNodeTable() {
	ctx := context.Background()
	host, _ := GetP2PNode(nil)
	// broadcast the peer to the network
	store, _ := GetCRDTStore()
	key := ds.NewKey(host.ID().String())
	common.Logger.Info("Removing myself from the network")
	store.Delete(ctx, key)
}

func UpdateNodeTableHook(key ds.Key, value []byte) {
	table := *GetNodeTable()
	var peer Peer
	err := json.Unmarshal(value, &peer)
	common.ReportError(err, "Error while unmarshalling peer")
	table[key.String()] = peer
}

func DeleteNodeTableHook(key ds.Key) {
	table := *GetNodeTable()
	delete(table, key.String())
}

func GetPeerFromTable(peerId string) (Peer, error) {
	table := *GetNodeTable()
	peer, ok := table[peerId]
	if !ok {
		return Peer{}, errors.New("peer not found")
	}
	return peer, nil
}

func GetService(name string) (Service, error) {
	host, _ := GetP2PNode(nil)
	store, _ := GetCRDTStore()
	key := ds.NewKey(host.ID().String())
	value, err := store.Get(context.Background(), key)
	common.ReportError(err, "Error while getting peer")
	var peer Peer
	err = json.Unmarshal(value, &peer)
	common.ReportError(err, "Error while unmarshalling peer")
	for _, service := range peer.Service {
		if service.Name == name {
			return service, nil
		}
	}
	return Service{}, errors.New("Service not found")
}

func GetAllProviders(serviceName string) ([]Peer, error) {
	var providers []Peer
	table := *GetNodeTable()
	for _, peer := range table {
		if peer.Connected {
			for _, service := range peer.Service {
				if service.Name == serviceName {
					providers = append(providers, peer)
				}
			}
		}
	}
	if len(providers) == 0 {
		return providers, errors.New("no providers found")
	}
	return providers, nil
}

func InitializeMyself() {
	host, _ := GetP2PNode(nil)
	ctx := context.Background()
	store, _ := GetCRDTStore()
	key := ds.NewKey(host.ID().String())
	myself = Peer{
		ID:            host.ID().String(),
		PublicAddress: viper.GetString("public-addr"),
	}
	myself.Hardware.GPUs = platform.GetGPUInfo()
	value, err := json.Marshal(myself)
	common.ReportError(err, "Error while marshalling peer")
	store.Put(ctx, key, value)
}
