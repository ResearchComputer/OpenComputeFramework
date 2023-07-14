package p2p

import (
	"encoding/json"
	"fmt"
	"ocfcore/internal/common"
	"ocfcore/internal/protocol/remote"

	"github.com/spf13/viper"
)

func BroadcastPeerOffering(peer Peer) {
	dnt := GetNodeTable()
	for _, remote := range dnt.Peers {
		if peer.PeerID != remote.PeerID {
			// we don't need to update local node table as it is already updated
			UpdateRemoteNodeTable(peer.PeerID, peer)
		}
	}
}

func UpdateRemoteNodeTable(peerId string, peer Peer) error {
	remoteAddr := fmt.Sprintf("http://localhost:%s/api/v1/proxy/%s/api/v1/status/table", viper.GetString("port"), peerId)
	reqString, err := json.Marshal(peer)
	if err != nil {
		return err
	}
	_, err = remote.HTTPPost(remoteAddr, reqString)
	if err != nil {
		common.Logger.Info("Error while updating remote node table", "error", err)
	}
	return err
}
