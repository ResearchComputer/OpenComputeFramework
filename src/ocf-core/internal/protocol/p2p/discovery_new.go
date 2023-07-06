package p2p

import (
	"context"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	routing "github.com/libp2p/go-libp2p/p2p/discovery/routing"
)

var discoverLockNew sync.Mutex

func DiscoverNew(ctx context.Context, h host.Host, dht *dht.IpfsDHT, rendezvous string) {
	discoverLockNew.Lock()
	defer discoverLockNew.Unlock()
	var routingDiscovery = routing.NewRoutingDiscovery(dht)
	routingDiscovery.Advertise(ctx, rendezvous)
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

}
