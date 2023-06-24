package queue

import (
	"errors"
	"ocfcore/internal/common"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

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
}

func Publish(topic string, data []byte) (*nats.Msg, error) {
	common.Logger.Info("Publishing to queue", "topic", topic, "data", string(data))
	msg, err := natsConn.Request(topic, data, 10*time.Second)
	return msg, err
}

func GetQueueStatus() (*server.Connz, error) {
	if natsServer == nil {
		common.Logger.Info("NATS server not started")
		return nil, errors.New("NATS server not started")
	}
	conn, err := natsServer.Connz(&server.ConnzOptions{Subscriptions: true, Offset: 1})
	return conn, err
}
