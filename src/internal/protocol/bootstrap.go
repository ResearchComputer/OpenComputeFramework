package protocol

import (
	"encoding/json"
	"ocf/internal/common"
	"strings"

	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/viper"
)

func getDefaultBootstrapPeers(bootstrapAddrs []string, mode string) []multiaddr.Multiaddr {
	var DefaultBootstrapPeers []multiaddr.Multiaddr
	if mode == "standalone" {
		bootstrapAddrs = []string{}
	} else if mode == "local" {
		bootstrapAddrs = []string{"/ip4/127.0.0.1/tcp/43905"}
	} else if bootstrapAddrs == nil {
		// read bootstrap_addr from config
		var bootstraps common.Bootstraps
		// send get request to this address
		resp, _ := common.RemoteGET(viper.GetString("bootstrap.addr"))
		json.Unmarshal(resp, &bootstraps)
		// if we got bootstrap addresses from remote, log it
		bootstrapAddrs = bootstraps.Bootstraps
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
