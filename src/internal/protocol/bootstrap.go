package protocol

import (
	"strings"

	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/viper"
)

const defaultBootstrapPeerAddr = "/ip4/140.238.223.13/tcp/43905/p2p/QmWxgDBrscNmiURmba196goATfG6fHrMniNDMei13YTCay"

func getDefaultBootstrapPeers(bootstrapAddrs []string, mode string) []multiaddr.Multiaddr {
	var DefaultBootstrapPeers []multiaddr.Multiaddr
	if mode == "standalone" {
		bootstrapAddrs = []string{}
	} else if mode == "local" {
		bootstrapAddrs = []string{"/ip4/127.0.0.1/tcp/43905"}
	} else if bootstrapAddrs == nil {
		// read bootstrap_addr from config
		configuredBootstrapAddrs := viper.GetStringSlice("bootstrap.addr")
		if len(configuredBootstrapAddrs) > 0 {
			bootstrapAddrs = append(bootstrapAddrs, configuredBootstrapAddrs...)
		} else {
			bootstrapAddrs = []string{defaultBootstrapPeerAddr}
		}
	}
	for _, s := range bootstrapAddrs {
		s = strings.TrimPrefix(s, "[")
		s = strings.TrimSuffix(s, "]")
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		DefaultBootstrapPeers = append(DefaultBootstrapPeers, ma)
	}
	return DefaultBootstrapPeers
}
