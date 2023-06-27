package server

import "ocfcore/internal/server/queue"

func DisconnectionDetection() {
	// List all connections
	queue.RemoveDisconnectedNode()
}
