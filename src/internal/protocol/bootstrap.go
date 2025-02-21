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
		bootstrapAddr := viper.GetString("bootstrap.addr")
		// if startswith http, send GET request
		if strings.HasPrefix(bootstrapAddr, "http") {
			common.Logger.Info("Sending GET request to: ", viper.GetString("bootstrap.addr"))
			resp, _ := common.RemoteGET(viper.GetString("bootstrap.addr"))
			err := json.Unmarshal(resp, &bootstraps)
			if err != nil {
				// convert resp to string
				common.Logger.Error("Failed to unmarshal bootstrap addresses: ", err)
				common.Logger.Info("Got response: ", string(resp))
			}
			// if we got bootstrap addresses from remote, log it
			bootstrapAddrs = bootstraps.Bootstraps
			common.Logger.Info("Got bootstrap addrs: ", bootstrapAddrs)
		} else if strings.HasPrefix(bootstrapAddr, "/ip4/") {
			bootstrapAddrs = []string{bootstrapAddr}
		} else {
			common.Logger.Error("Invalid bootstrap address: ", bootstrapAddr)
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
	common.Logger.Info("Bootstrap: ", DefaultBootstrapPeers)
	return DefaultBootstrapPeers
}
