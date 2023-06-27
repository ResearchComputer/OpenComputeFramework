package requests

import (
	"ocfcore/internal/common/structs"
	"ocfcore/internal/server/p2p"
)

// functions for massively broadcasting messages to all peers

func BroadcastNodeStatus(nodeStatus structs.NodeStatus) {
	node := p2p.GetP2PNode()
	peers := node.Peerstore().Peers()
	for _, peer := range peers {
		if peer.String() != node.ID().String() {
			// we don't need to update local node table as it is already updated
			UpdateRemoteNodeTable(peer.String(), nodeStatus)
		}
	}
}
