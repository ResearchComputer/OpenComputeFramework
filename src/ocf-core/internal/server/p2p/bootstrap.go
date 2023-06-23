package p2p

import "github.com/multiformats/go-multiaddr"

func getDefaultBootstrapPeers() []multiaddr.Multiaddr {
	var DefaultBootstrapPeers []multiaddr.Multiaddr
	for _, s := range []string{
		"/ip4/164.90.231.186/tcp/43905/p2p/QmbY2bk4JGkD6yoW9DriYsFqHqqSjZh7AyyuXeYYKFDXba",
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		DefaultBootstrapPeers = append(DefaultBootstrapPeers, ma)
	}
	return DefaultBootstrapPeers
}
