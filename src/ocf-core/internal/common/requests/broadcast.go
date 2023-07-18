package requests

import (
	"ocfcore/internal/common/structs"
	"ocfcore/internal/protocol/p2p"
)

// functions for massively broadcasting messages to all peers

func BroadcastNodeStatus(nodeStatus structs.NodeStatus) {
	node := p2p.GetP2PNode()
	dnt := p2p.GetNodeTable()
	for _, peer := range dnt.Peers {
		if peer.PeerID != node.ID().String() {
			// we don't need to update local node table as it is already updated
			UpdateRemoteNodeTable(peer.PeerID, nodeStatus)
		}
	}
}
