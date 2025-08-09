package protocol

import (
	"encoding/json"
	"math/rand"
	"ocf/internal/common"
	"ocf/internal/common/process"
	"os"

	ds "github.com/ipfs/go-datastore"
	"github.com/jasonlvhit/gocron"
	"github.com/libp2p/go-libp2p/core/network"
)

// var verificationKey = "ocf-verification-key"
var verificationProb = 0.5

func StartTicker() {
	err := gocron.Every(1).Minute().Do(func() {
		common.Logger.Info("Starting verification")
		if rand.Float64() < verificationProb {
			// store, _ := GetCRDTStore()
			// ctx := context.Background()
			// store.Put(ctx, ds.NewKey(verificationKey), []byte("verification"))
			Reconnect()
		}
	})
	common.ReportError(err, "Error while creating verification ticker")
	err = gocron.Every(30).Second().Do(func() {
		host, _ := GetP2PNode(nil)
		peers := host.Peerstore().Peers()
		// updateMyself()
		for _, peer_id := range peers {
			// check if peer is still connected
			peer, error := GetPeerFromTable(peer_id.String())
			if error == nil {
				peer.Connected = true
				if peer_id != host.ID() && host.Network().Connectedness(peer_id) != network.Connected {
					common.Logger.Info("Peer:" + peer_id.String() + " got disconnected!")
					peer.Connected = false
				}
				value, err := json.Marshal(peer)
				if err == nil {
					UpdateNodeTableHook(ds.NewKey(peer_id.String()), value)
				}
			}
		}
		if !process.HealthCheck() {
			common.Logger.Error("Health check failed")
			// exit myself
			os.Exit(1)
		}
	})
	common.ReportError(err, "Error while creating cleaning ticker")
	<-gocron.Start()
}
