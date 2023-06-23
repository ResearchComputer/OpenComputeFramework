package queue

import (
	"ocfcore/internal/common"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

var natsConn *nats.Conn

func StartQueueServer() {
	common.Logger.Info("Starting queue server...")
	opts := &server.Options{
		JetStream: true,
		Port:      viper.GetInt("queue.port"),
	}
	ns, err := server.NewServer(opts)

	if err != nil {
		panic(err)
	}
	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}
	natsConn, err = nats.Connect(ns.ClientURL())
	if err != nil {
		panic(err)
	}
}

func Publish(topic string, data []byte) (*nats.Msg, error) {
	common.Logger.Info("Publishing to queue", "topic", topic, "data", string(data))
	msg, err := natsConn.Request(topic, data, 10*time.Second)
	return msg, err
}
