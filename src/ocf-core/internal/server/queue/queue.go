package queue

import (
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

func StartQueueServer() {
	opts := &server.Options{}
	ns, err := server.NewServer(opts)

	if err != nil {
		panic(err)
	}
	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}
}
