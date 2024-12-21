package server

import (
	"ocf/internal/protocol"
	"time"

	"github.com/axiomhq/axiom-go/axiom"
	"github.com/axiomhq/axiom-go/axiom/ingest"
	"github.com/gin-gonic/gin"
)

func listPeers(c *gin.Context) {
	addrs := protocol.ConnectedPeers()
	c.JSON(200, gin.H{"peers": addrs})
}

func listBootstraps(c *gin.Context) {
	addrs := protocol.ConnectedBootstraps()
	c.JSON(200, gin.H{"bootstraps": addrs})
}

func updateLocal(c *gin.Context) {
	var peer protocol.Peer
	c.BindJSON(&peer)
	protocol.UpdateNodeTable(peer)
}

func deleteLocal(c *gin.Context) {
	var peer protocol.Peer
	c.BindJSON(&peer)
	protocol.DeleteNodeTable()
}

func getDNT(c *gin.Context) {
	events := []axiom.Event{
		{ingest.TimestampField: time.Now(), "event": "DNT Lookup"},
	}
	IngestEvents(events)
	c.JSON(200, protocol.GetNodeTable(true))
}
