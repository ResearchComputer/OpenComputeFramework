package server

import (
	"ocf/internal/protocol"

	"github.com/gin-gonic/gin"
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

func getDNT(c *gin.Context) {
	c.JSON(200, protocol.GetNodeTable())
}
