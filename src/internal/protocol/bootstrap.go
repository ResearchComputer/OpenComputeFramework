package protocol

import (
	"github.com/multiformats/go-multiaddr"
)

const defaultBootstrapPeerAddr = "/ip4/206.189.249.2/tcp/43905/p2p/QmbY2bk4JGkD6yoW9DriYsFqHqqSjZh7AyyuXeYYKFDXba"

// GetDefaultBootstrapPeers returns the default bootstrap peers.
func getDefaultBootstrapPeers(bootstrapAddrs []string, mode string) []multiaddr.Multiaddr {
	var DefaultBootstrapPeers []multiaddr.Multiaddr
	if mode == "standalone" {
		bootstrapAddrs = []string{}
	} else if bootstrapAddrs == nil {
		bootstrapAddrs = []string{defaultBootstrapPeerAddr}
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
