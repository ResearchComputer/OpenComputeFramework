package p2p

import (
	"context"
	"net"
	"ocfcore/internal/common"

	gostream "github.com/libp2p/go-libp2p-gostream"
	p2phttp "github.com/libp2p/go-libp2p-http"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/viper"
)

func P2PListener() net.Listener {
	ctx := context.Background()
	host := GetP2PNode()
	var dhtc *dht.IpfsDHT
	var err error
	if viper.GetString("bootstrap.mode") == "standalone" {
		common.Logger.Info("standalone mode")
		dhtc, err = NewDHT(ctx, host, []multiaddr.Multiaddr{})
	} else {
		dhtc, err = NewDHT(ctx, host, getDefaultBootstrapPeers())
	}
	if err != nil {
		panic(err)
	}
	dhtc.Bootstrap(ctx)
	go Discover(ctx, host, dhtc, viper.GetString("bootstrap.rendezvous"))
	common.Logger.Info("ocfcore peer ID: ", host.ID())
	common.Logger.Info("ocfcore peer Addr: ", host.Addrs())
	listener, _ := gostream.Listen(host, p2phttp.DefaultP2PProtocol)
	return listener
}
