package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"ocfcore/internal/common"
	"ocfcore/internal/common/structs"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

// todo(xiaozhe): this probably needs a refactor with singleton pattern

var natsConn *nats.Conn
var natsServer *server.Server

func StartQueueServer() {
	var err error
	common.Logger.Info("Starting queue server...")
	opts := &server.Options{
		JetStream: true,
		Port:      viper.GetInt("queue.port"),
	}
	natsServer, err = server.NewServer(opts)
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
	common.Logger.Debug("Publishing to queue", "topic", topic, "data", string(data))
	msg, err := natsConn.Request(topic, data, 10*time.Second)
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
	_, err := natsConn.Subscribe("worker:status", func(msg *nats.Msg) {
		var nodeStatus structs.NodeStatus
		err := json.Unmarshal(msg.Data, &nodeStatus)
		if err != nil {
			common.Logger.Error("Failed to unmarshal worker status", "error", err)
		}
		fmt.Println(nodeStatus)
	})
	return err
}

func GetProvidedService() ([]string, error) {
	if natsServer == nil {
		common.Logger.Info("NATS server not started")
		return nil, errors.New("NATS server not started")
	}
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
