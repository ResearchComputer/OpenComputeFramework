package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"ocfcore/internal/common"
	"ocfcore/internal/common/structs"
	"ocfcore/internal/protocol/p2p"
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
var lnt *p2p.NodeTable

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

func NewNodeTable() *p2p.NodeTable {
	tableOnce.Do(func() {
		fmt.Println("Initializing node table")
		lnt = &p2p.NodeTable{}
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
		common.Logger.Debug("Received worker status", "data", string(msg.Data))
		var nodeStatus structs.NodeStatus
		err := json.Unmarshal(msg.Data, &nodeStatus)
		if err != nil {
			common.Logger.Error("Failed to unmarshal worker status", "error", err)
		}
		nodeStatus.PeerID = p2p.GetP2PNode().ID().String()
		if nodeStatus.Status == "connected" {
			hw := p2p.HardwareSpec{
				GPUs: nodeStatus.Specs,
			}
			p2p.GetNodeTable().NewOffering(nodeStatus.PeerID, nodeStatus.Service, hw)
		} else if nodeStatus.Status == "disconnected" {
			p2p.GetNodeTable().RemoveOffering(nodeStatus.PeerID, nodeStatus.Service)
		}
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
