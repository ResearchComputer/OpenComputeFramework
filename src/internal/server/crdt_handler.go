package server

import (
	"ocf/internal/protocol"

	"github.com/gin-gonic/gin"
)

func listPeers(c *gin.Context) {
	// crdtNode, cancelCtx := protocol.GetCRDTStore()
	// defer cancelCtx()
	addrs := protocol.ConnectedPeers()
	c.JSON(200, gin.H{"peers": addrs})
}
