package p2p

import (
	"context"
	"rccore/internal/common"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func NewDHT(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr) (*dht.IpfsDHT, error) {
	var options []dht.Option

	if len(bootstrapPeers) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
	}

	kdht, err := dht.New(ctx, host, options...)
	if err != nil {
		return nil, err
	}

	if err = kdht.Bootstrap(ctx); err != nil {
		return nil, err
	}

	for _, peerAddr := range bootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		go func() {
			if err := host.Connect(ctx, *peerinfo); err != nil {
				common.Logger.Warn("Error while connecting to node %q: %-v", peerinfo, err)
			}
		}()
	}

	return kdht, nil
}
