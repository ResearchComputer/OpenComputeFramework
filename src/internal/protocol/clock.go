package protocol

import (
	"math/rand"
	"ocf/internal/common"

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
		for _, peer := range peers {
			// check if peer is still connected
			if host.Network().Connectedness(peer) != network.Connected {
				// delete peer from table
				DeleteNodeTableHook(ds.NewKey(peer.String()))
			}
		}
	})
	common.ReportError(err, "Error while creating cleaning ticker")
	<-gocron.Start()
}
