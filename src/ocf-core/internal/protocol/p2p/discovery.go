package p2p

import (
	"context"
	"ocfcore/internal/common"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peerstore"
	routing "github.com/libp2p/go-libp2p/p2p/discovery/routing"
)

var DisconnectedPeers []string
var discoverLock sync.Mutex

func Discover(ctx context.Context, h host.Host, dht *dht.IpfsDHT, rendezvous string) {
	discoverLock.Lock()
	defer discoverLock.Unlock()
	var routingDiscovery = routing.NewRoutingDiscovery(dht)

	routingDiscovery.Advertise(ctx, rendezvous)

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// cleaning peerstore first
			storedPeers := h.Peerstore().Peers()
			for _, p := range storedPeers {
				if h.Network().Connectedness(p) == network.NotConnected {
					DisconnectedPeers = append(DisconnectedPeers, p.String())
				}
			}
			peers, err := routingDiscovery.FindPeers(ctx, rendezvous)
			if err != nil {
				common.Logger.Error(err)
			}
			for p := range peers {
				if p.ID == h.ID() {
					continue
				}
				if h.Network().Connectedness(p.ID) != network.Connected {
					_, err = h.Network().DialPeer(ctx, p.ID)
					if err != nil {
						continue
					}
				}
				common.Logger.Debug("Connectivity to peer: ", p.ID, " is ", h.Network().Connectedness(p.ID))
				if h.Network().Connectedness(p.ID) == network.Connected {
					h.Peerstore().AddAddrs(p.ID, p.Addrs, peerstore.PermanentAddrTTL)
					// remove it from DisconnectedPeers
					for i, dp := range DisconnectedPeers {
						if dp == p.ID.String() {
							DisconnectedPeers = append(DisconnectedPeers[:i], DisconnectedPeers[i+1:]...)
							break
						}
					}
				}
			}
		}
	}
}
