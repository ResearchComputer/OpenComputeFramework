package requests

import (
	"fmt"
	"io"
	"log"

	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var blackListedPeers []string

func CheckPeerStatus(peerId string) error {
	if !slices.Contains(blackListedPeers, peerId) {
		peerAddr := fmt.Sprintf("http://localhost:%s/api/v1/proxy/%s/api/v1/status/health?peer=0", viper.GetString("port"), peerId)
		resp, err := NewHTTPClient().Get(peerAddr)
		if err != nil {
			log.Fatalln(err)
			return err
		} else {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
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
	}
	return fmt.Errorf("peer is blacklisted")
}
