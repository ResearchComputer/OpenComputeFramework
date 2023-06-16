package p2p

import (
	"context"
	"crypto/rand"
	"io"
	mrand "math/rand"
	"ocfcore/internal/common"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
)

var P2PNode host.Host

func GetP2PNode() host.Host {
	if P2PNode == nil {
		ctx := context.Background()
		var err error
		P2PNode, err = newHost(ctx, 0)
		if err != nil {
			panic(err)
		}
	}
	return P2PNode
}

func newHost(ctx context.Context, seed int64) (host.Host, error) {
	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		common.Logger.Error("Error while creating connection manager: %v", err)
	}

	// try to load the private key from file
	priv := loadKeyFromFile()
	if priv == nil {
		// if it doesn't exist, generate a new one
		// If the seed is zero, use real cryptographic randomness. Otherwise, use a
		// deterministic randomness source to make generated keys stay the same
		// across multiple runs
		var r io.Reader
		if seed == 0 {
			r = rand.Reader
		} else {
			r = mrand.New(mrand.NewSource(seed))
		}
		priv, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
		if err != nil {
			return nil, err
		}
	}

	// persist private key
	writeKeyToFile(priv)
	if err != nil {
		return nil, err
	}

	return libp2p.New(
		libp2p.DefaultTransports,
		libp2p.Identity(priv),
		libp2p.ConnectionManager(connmgr),
		libp2p.NATPortMap(),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/43905",
			"/ip4/0.0.0.0/udp/59820/quic",
		),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(noise.ID, noise.New),
		libp2p.EnableNATService(),
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
		libp2p.ForceReachabilityPublic(),
	)
}
