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
	common.Logger.Info("Broadcasting peer offering", "peer", peer)
	for _, remote := range dnt.Peers {
		if peer.PeerID != remote.PeerID {
			UpdateRemoteNodeTable(remote.PeerID, peer)
		}
	}
}

func UpdateRemoteNodeTable(peerId string, peer Peer) error {
	peer.Owner = viper.GetString("wallet.account")
	remoteAddr := fmt.Sprintf("http://localhost:%s/api/v1/proxy/%s/api/v1/status/peers", viper.GetString("port"), peerId)
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
