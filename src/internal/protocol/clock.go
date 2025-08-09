package protocol

import (
	"encoding/json"
	"math/rand"
	"ocf/internal/common"
	"ocf/internal/common/process"
	"os"
	"time"

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
				// update last seen timestamp
				peer.LastSeen = time.Now().Unix()
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
			Reconnect()
			// best-effort re-announce our services after trying to reconnect
			ReannounceLocalServices()
		}

		// Cleanup: remove peers that have been disconnected for a long time
		// Define staleness threshold
		staleAfter := 10 * time.Minute
		table := *GetNodeTable()
		now := time.Now().Unix()
		for id, p := range table {
			if !p.Connected && p.LastSeen > 0 {
				if time.Unix(p.LastSeen, 0).Add(staleAfter).Before(time.Now()) {
					common.Logger.Warnf("Removing stale peer %s (last seen %v)", id, time.Unix(p.LastSeen, 0))
					DeleteNodeTableHook(ds.NewKey(id))
				}
			}
			// Also mark peers with very old LastSeen as disconnected
			if p.Connected && p.LastSeen > 0 && time.Unix(p.LastSeen, 0).Add(2*time.Minute).Before(time.Now()) {
				p.Connected = false
				value, err := json.Marshal(p)
				if err == nil {
					UpdateNodeTableHook(ds.NewKey(id), value)
				}
			}
			// If LastSeen is zero, initialize it now
			if p.LastSeen == 0 {
				p.LastSeen = now
				value, err := json.Marshal(p)
				if err == nil {
					UpdateNodeTableHook(ds.NewKey(id), value)
				}
			}
		}
	})
	common.ReportError(err, "Error while creating resource monitoring ticker")
	<-gocron.Start()
}
