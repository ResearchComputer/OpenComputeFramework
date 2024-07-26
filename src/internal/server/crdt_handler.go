package server

import (
	"ocf/internal/protocol"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func listPeers(c *gin.Context) {
	addrs := protocol.ConnectedPeers()
	c.JSON(200, gin.H{"peers": addrs})
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
	_, span := tracer.Start(c.Request.Context(), "getDNT", oteltrace.WithAttributes(attribute.String("id", "test")))
	defer span.End()
	c.JSON(200, protocol.GetNodeTable())
}
