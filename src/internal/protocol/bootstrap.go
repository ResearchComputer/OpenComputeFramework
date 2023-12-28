package protocol

import (
	"github.com/multiformats/go-multiaddr"
)

const defaultBootstrapPeerAddr = "/ip4/140.238.223.13/tcp/43905/p2p/QmWxgDBrscNmiURmba196goATfG6fHrMniNDMei13YTCay"

// GetDefaultBootstrapPeers returns the default bootstrap peers.
func getDefaultBootstrapPeers(bootstrapAddrs []string, mode string) []multiaddr.Multiaddr {
	var DefaultBootstrapPeers []multiaddr.Multiaddr
	if mode == "standalone" {
		bootstrapAddrs = []string{}
	} else if mode == "local" {
		bootstrapAddrs = []string{"/ip4/127.0.0.1/tcp/43905"}
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
