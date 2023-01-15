package p2p

import "github.com/multiformats/go-multiaddr"

func getDefaultBootstrapPeers() []multiaddr.Multiaddr {
	var DefaultBootstrapPeers []multiaddr.Multiaddr
	for _, s := range []string{
		"/ip4/140.238.214.135/tcp/43905/p2p/QmbQA8PZC5hFguetaUfGWuJyfqCKnRFPuNH6htDthiAzUs",
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		DefaultBootstrapPeers = append(DefaultBootstrapPeers, ma)
	}
	return DefaultBootstrapPeers
}
