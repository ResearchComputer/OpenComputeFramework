package protocol

import (
	"github.com/multiformats/go-multiaddr"
)

const defaultBootstrapPeerAddr = "/ip4/206.189.249.2/tcp/43905/p2p/QmbY2bk4JGkD6yoW9DriYsFqHqqSjZh7AyyuXeYYKFDXba"

// GetDefaultBootstrapPeers returns the default bootstrap peers.
func getDefaultBootstrapPeers(bootstrapAddrs []string) []multiaddr.Multiaddr {
	var DefaultBootstrapPeers []multiaddr.Multiaddr
	if bootstrapAddrs == nil {
		bootstrapAddrs = append(bootstrapAddrs, defaultBootstrapPeerAddr)
	}
	for _, s := range bootstrapAddrs {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		DefaultBootstrapPeers = append(DefaultBootstrapPeers, ma)
	}
	return DefaultBootstrapPeers
}
