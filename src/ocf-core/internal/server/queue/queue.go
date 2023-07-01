package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"ocfcore/internal/common"
	"ocfcore/internal/common/requests"
	"ocfcore/internal/common/structs"
	"ocfcore/internal/server/p2p"
	"sync"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

// package-level variables
// read-only by external packages through New* functions
// read-write by current package
var natsConn *nats.Conn
var natsServer *server.Server
var lnt *structs.NodeTable

var natsOnce sync.Once
var tableOnce sync.Once

func NewNatsServer() *server.Server {
	natsOnce.Do(func() {
		var err error
		opts := &server.Options{
			JetStream: true,
			Port:      viper.GetInt("queue.port"),
		}
		natsServer, err = server.NewServer(opts)
		if err != nil {
			panic(err)
		}
	})
	return natsServer
}

func NewNodeTable() *structs.NodeTable {
	tableOnce.Do(func() {
		fmt.Println("Initializing node table")
		lnt = &structs.NodeTable{}
	})
	return lnt
}

func StartQueueServer() {
	var err error
	common.Logger.Info("Starting queue server...")
	natsServer = NewNatsServer()
	if err != nil {
		panic(err)
	}
	natsServer.Start()
	if !natsServer.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}
	natsConn, err = nats.Connect(natsServer.ClientURL())
	if err != nil {
		panic(err)
	}
	common.Logger.Info("Queue server listening on port: ", viper.GetInt("queue.port"))
	SubscribeWorkerStatus()
}

func Publish(topic string, data []byte) (*nats.Msg, error) {
	msg, err := natsConn.Request(topic, data, 3600*time.Second)
	return msg, err
}

func GetQueueStatus() (*server.Connz, error) {
	if natsServer == nil {
		common.Logger.Debug("NATS server not started")
		return nil, errors.New("NATS server not started")
	}
	conn, err := natsServer.Connz(&server.ConnzOptions{Subscriptions: true, Offset: 1})
	return conn, err
}

func SubscribeWorkerStatus() error {
	if natsConn == nil {
		common.Logger.Debug("NATS client not started")
		return nil
	}
	var lock = &sync.Mutex{}
	_, err := natsConn.Subscribe("worker:status", func(msg *nats.Msg) {
		lock.Lock()
		// make sure the table is being updated by only one worker at a time
		defer lock.Unlock()

		var nodeStatus structs.NodeStatus
		err := json.Unmarshal(msg.Data, &nodeStatus)
		if err != nil {
			common.Logger.Error("Failed to unmarshal worker status", "error", err)
		}
		nodeStatus.PeerID = p2p.GetP2PNode().ID().String()
		table := NewNodeTable()
		*lnt = *table.Update(nodeStatus)
		go requests.BroadcastNodeStatus(nodeStatus)
	})
	return err
}

func GetProvidedService() ([]string, error) {
	natsServer := NewNatsServer()
	conn, err := natsServer.Connz(&server.ConnzOptions{Subscriptions: true, Offset: 1})
	if err != nil {
		return nil, err
	}
	var providedService []string
	for _, c := range conn.Conns {
		providedService = append(providedService, c.Subs...)
	}
	return providedService, nil
}

func UpdateNodeTable(nodeStatus structs.NodeStatus) {
	*lnt = *NewNodeTable().Update(nodeStatus)
}

func RemovePeerFromNodeTable(peerID string) {
	for _, node := range NewNodeTable().Nodes {
		if node.PeerID == peerID {
			node.Status = "disconnected"
			*lnt = *NewNodeTable().Update(node)
		}
	}
}

func RemoveDisconnectedNode() {
	natsServer := NewNatsServer()
	p2pNode := p2p.GetP2PNode()
	conn, err := natsServer.Connz(&server.ConnzOptions{Subscriptions: true, Offset: 1})
	if err != nil {
		common.Logger.Error("Failed to get connection status: ", err)
		common.Logger.Error("If this persists, the node table will not be updated")
		return
	}
	// two steps:
	// check if all nodes in the DNT are still connected
	// for all nodes in lnt, check if my workers are still connected
	for _, node := range NewNodeTable().Nodes {
		// if it is the current node, then continue
		if node.PeerID == p2pNode.ID().String() {
			connected := false
			for _, c := range conn.Conns {
				if c.Cid == uint64(node.ClientID) {
					connected = true
					break
				}
			}
			if !connected {
				// if not connected, remove from node table
				node.Status = "disconnected"
				*lnt = *NewNodeTable().Update(node)
				go requests.BroadcastNodeStatus(node)
			}
		} else {
			// check if it is in peerstore
			disconnected := false
			for _, p := range p2p.DisconnectedPeers {
				if p == node.PeerID {
					disconnected = true
					break
				}
			}
			if disconnected {
				common.Logger.Debug("Peer ", node.PeerID, " is disconnected")
				// if disconnected, remove from node table
				node.Status = "disconnected"
				*lnt = *NewNodeTable().Update(node)
			}
		}

	}
}
