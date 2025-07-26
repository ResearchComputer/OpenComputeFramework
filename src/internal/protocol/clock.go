package protocol

import (
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
		for _, peer := range peers {
			// check if peer is still connected
			if peer != host.ID() && host.Network().Connectedness(peer) != network.Connected {
				// delete peer from table
				DeleteNodeTableHook(ds.NewKey(peer.String()))
			}
		}
		if !process.HealthCheck() {
			common.Logger.Error("Health check failed")
			// exit myself
			os.Exit(1)
		}
	})
	common.ReportError(err, "Error while creating cleaning ticker")

	// Add resource monitoring every 2 minutes
	err = gocron.Every(2).Minutes().Do(func() {
		GetResourceManagerStats()

		// Also log current connection count for easy monitoring
		connectedPeers := ConnectedPeers()
		allPeers := AllPeers()
		common.Logger.Infof("Connection Summary: %d connected peers, %d total known peers",
			len(connectedPeers), len(allPeers))

		// Log if we have very few connections (potential issue)
		if len(connectedPeers) < 3 {
			common.Logger.Warnf("Low connection count detected: only %d connected peers", len(connectedPeers))
		}
	})
	common.ReportError(err, "Error while creating resource monitoring ticker")
	<-gocron.Start()
}
