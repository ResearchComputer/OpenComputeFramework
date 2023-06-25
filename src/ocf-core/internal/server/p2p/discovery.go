package p2p

import (
	"context"
	"ocfcore/internal/common"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peerstore"
	routing "github.com/libp2p/go-libp2p/p2p/discovery/routing"
)

func Discover(ctx context.Context, h host.Host, dht *dht.IpfsDHT, rendezvous string) {
	var routingDiscovery = routing.NewRoutingDiscovery(dht)

	routingDiscovery.Advertise(ctx, rendezvous)

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
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
				if h.Network().Connectedness(p.ID) == network.Connected {
					common.Logger.Info("Connected to: ", p.ID)
					h.Peerstore().AddAddrs(p.ID, p.Addrs, peerstore.PermanentAddrTTL)
				}
			}
		}
	}
}
