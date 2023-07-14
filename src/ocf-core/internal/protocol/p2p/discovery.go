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

var discoverLock sync.Mutex

// DiscoverNew is a function that keeps updating DNT with the latest information about the network.
func Discover(ctx context.Context, h host.Host, dht *dht.IpfsDHT, rendezvous string) {
	GetNodeTable().Update(Peer{
		PeerID: GetP2PNode().ID().String(),
		Status: CONNECTED,
	})
	var disconnected []string
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
			// cleaning disconnected peers
			storedPeers := h.Peerstore().Peers()
			for _, p := range storedPeers {
				if p != h.ID() && h.Network().Connectedness(p) == network.NotConnected {
					disconnected = append(disconnected, p.String())
				}
			}
			GetNodeTable().RemoveDisconnectedPeers(disconnected)
			disconnected = []string{}
			peers, err := routingDiscovery.FindPeers(ctx, rendezvous)
			if err != nil {
				common.Logger.Error(err)
			}
			for p := range peers {
				if p.ID == h.ID() {
					continue
				}
				if h.Network().Connectedness(p.ID) != network.Connected {
					_, err := h.Network().DialPeer(ctx, p.ID)
					if err != nil {
						continue
					}
				}
				if h.Network().Connectedness(p.ID) == network.Connected {
					h.Peerstore().AddAddrs(p.ID, p.Addrs, peerstore.PermanentAddrTTL)
					GetNodeTable().Update(Peer{
						PeerID: p.ID.String(),
						Status: CONNECTED,
					})
				}
			}
		}
	}
}
