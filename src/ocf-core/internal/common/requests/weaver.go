package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"ocfcore/internal/common"

	ns "github.com/nats-io/nats-server/v2/server"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var blackListedPeers []string

func CheckPeerStatus(peerId string) error {
	if !slices.Contains(blackListedPeers, peerId) {
		peerAddr := fmt.Sprintf("http://localhost:%s/api/v1/proxy/%s/api/v1/status/health?peer=0", viper.GetString("port"), peerId)
		resp, err := NewHTTPClient().Get(peerAddr)
		if err != nil {
			common.Logger.Error("Error while checking peer status", "error", err)
			return err
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			common.Logger.Error("Error while reading response body", "error", err)
			return err
		}
		if string(b) == "ERROR: protocol not supported" {
			blackListedPeers = append(blackListedPeers, peerId)
			return fmt.Errorf("peer %s is not ocfcore", peerId)
		} else {
			fmt.Println(string(b))
		}
		return nil

	}
	return fmt.Errorf("peer is blacklisted")
}

func ReadProvidedService(peerId string) ([]string, error) {
	remoteAddr := fmt.Sprintf("http://localhost:%s/api/v1/proxy/%s/api/v1/status/connections", viper.GetString("port"), peerId)
	resp, err := NewHTTPClient().Get(remoteAddr)
	if err != nil {
		common.Logger.Debug("Error while reading provided service", "error", err)
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		common.Logger.Error("Error while reading response body", "error", err)
		return nil, err
	}
	fmt.Println(string(b))
	var conns ns.Connz
	json.Unmarshal(b, &conns)
	var providedService []string
	for _, conn := range conns.Conns {
		providedService = append(providedService, conn.Subs...)
	}
	return providedService, nil
}
