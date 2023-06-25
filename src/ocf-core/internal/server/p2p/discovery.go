package p2p

import (
	"context"
	"ocfcore/internal/common"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
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
				common.Logger.Debug("Found peer: ", p)
				if p.ID == h.ID() {
					continue
				}
				common.Logger.Debug("Peer connectedness: ", h.Network().Connectedness(p.ID))
				if h.Network().Connectedness(p.ID) != network.Connected {
					_, err = h.Network().DialPeer(ctx, p.ID)
					common.Logger.Info("Dialing to: ", p)
					common.Logger.Error(err)
					if err != nil {
						continue
					}
					// add to peerstore
					common.Logger.Info("Connected to: ", p)
					h.Peerstore().AddAddrs(p.ID, p.Addrs, time.Hour)
				}
			}
		}
	}
}
