package protocol

import (
	"context"
	"crypto/rand"
	"encoding/json"
	mrand "math/rand"
	"ocf/internal/common"
	"strconv"
	"sync"
	"time"

	"github.com/ipfs/boxo/ipns"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dualdht "github.com/libp2p/go-libp2p-kad-dht/dual"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	connmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/spf13/viper"
)

var P2PNode *host.Host
var ddht *dualdht.DHT
var hostOnce sync.Once
var MyID string

func GetP2PNode(ds datastore.Batching) (host.Host, dualdht.DHT) {
	hostOnce.Do(func() {
		ctx := context.Background()
		var err error
		seed := viper.GetString("seed")
		// try to parse the seed as int64
		seedInt, err := strconv.ParseInt(seed, 10, 64)
		if err != nil {
			panic(err)
		}
		host, err := newHost(ctx, seedInt, ds)
		MyID = host.ID().String()
		P2PNode = &host
		if err != nil {
			panic(err)
		}
	})
	return *P2PNode, *ddht
}

func newHost(ctx context.Context, seed int64, ds datastore.Batching) (host.Host, error) {
	var err error
	// Connection manager: maintain a larger pool of connections so we can exceed
	// the pubsub mesh degree and keep more peers around.
	cm, err := connmgr.NewConnManager(
		100, // Low watermark
		800, // High watermark
		connmgr.WithGracePeriod(5*time.Minute),
	)
	if err != nil {
		common.Logger.Error("Error while creating connection manager: ", err)
	}
	var priv crypto.PrivKey
	// try to load the private key from file
	if seed == 0 {
		// try to load from the file
		priv = loadKeyFromFile()
		if priv == nil {
			r := rand.Reader
			priv, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
			if err != nil {
				return nil, err
			}
			writeKeyToFile(priv)
		}
	} else {
		r := mrand.New(mrand.NewSource(seed))
		priv, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
		if err != nil {
			return nil, err
		}
		writeKeyToFile(priv)
	}
	if err != nil {
		return nil, err
	}
	// Configure resource manager with higher limits
	limits := rcmgr.DefaultLimits.AutoScale()
	// Increase connection limits significantly for a distributed system
	systemLimits := rcmgr.ResourceLimits{
		ConnsInbound:    1000,  // Allow up to 1000 inbound connections
		ConnsOutbound:   1000,  // Allow up to 1000 outbound connections
		Conns:           2000,  // Allow up to 2000 total connections
		StreamsInbound:  10000, // Increase stream limits
		StreamsOutbound: 10000,
		Streams:         20000,
		Memory:          16 << 30, // 16GB memory limit
	}

	// Apply the custom limits
	finalLimits := rcmgr.PartialLimitConfig{
		System: systemLimits,
		// Keep default peer limits but increase them slightly
		PeerDefault: rcmgr.ResourceLimits{
			ConnsInbound:  512, // Allow more connections per peer
			ConnsOutbound: 512,
			Conns:         1024,
		},
	}.Build(limits)

	// Create resource manager
	mgr, err := rcmgr.NewResourceManager(rcmgr.NewFixedLimiter(finalLimits))
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.DefaultTransports,
		libp2p.Identity(priv),
		libp2p.ResourceManager(mgr), // Use our custom resource manager
		libp2p.ConnectionManager(cm),
		libp2p.NATPortMap(),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/"+viper.GetString("tcpport"),
			"/ip4/0.0.0.0/tcp/"+viper.GetString("tcpport")+"/ws",
			"/ip4/0.0.0.0/udp/"+viper.GetString("udpport")+"/quic",
		),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(noise.ID, noise.New),
		libp2p.EnableNATService(),
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
		libp2p.ForceReachabilityPublic(),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			ddht, err = newDHT(ctx, h, ds)
			return ddht, err
		}),
	}

	host, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	// Log connection events for debugging
	host.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(n network.Network, c network.Conn) {
			common.Logger.Info("Connected to peer: ", c.RemotePeer(), " Total connections: ", len(n.Conns()))
			// On (re)connections, re-announce local services
			go ReannounceLocalServices()

			// Mark peer as connected in node table immediately
			go func(pid peer.ID) {
				// Avoid updating self
				if pid == host.ID() {
					return
				}
				p, err := GetPeerFromTable(pid.String())
				if err != nil {
					p = Peer{ID: pid.String()}
				}
				p.Connected = true
				p.LastSeen = time.Now().Unix()
				if b, e := json.Marshal(p); e == nil {
					UpdateNodeTableHook(datastore.NewKey(pid.String()), b)
				} else {
					common.Logger.Error("Failed to marshal peer on connect: ", e)
				}
			}(c.RemotePeer())
		},
		DisconnectedF: func(n network.Network, c network.Conn) {
			common.Logger.Info("Disconnected from peer: ", c.RemotePeer(), " Total connections: ", len(n.Conns()))
			// Mark peer as disconnected in node table immediately
			go func(pid peer.ID) {
				if pid == host.ID() {
					return
				}
				p, err := GetPeerFromTable(pid.String())
				if err != nil {
					p = Peer{ID: pid.String()}
				}
				p.Connected = false
				// keep LastSeen as last known good; do not bump here
				if b, e := json.Marshal(p); e == nil {
					UpdateNodeTableHook(datastore.NewKey(pid.String()), b)
				} else {
					common.Logger.Error("Failed to marshal peer on disconnect: ", e)
				}
			}(c.RemotePeer())
		},
	})

	// Start a background auto-reconnector that watches connectivity
	go startAutoReconnect(ctx, host)

	return host, nil
}

// startAutoReconnect periodically checks if we lost connectivity and attempts to reconnect to bootstraps with backoff.
func startAutoReconnect(ctx context.Context, h host.Host) {
	// exponential backoff parameters
	minDelay := 5 * time.Second
	maxDelay := 2 * time.Minute
	delay := minDelay
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(30 * time.Second):
			// If very few or zero peers, try bootstrap
			conns := h.Network().Conns()
			if len(conns) == 0 {
				common.Logger.Warn("No active P2P connections; attempting reconnect to bootstraps...")
				Reconnect()
				// after a reconnect attempt, wait with backoff if still disconnected
				time.Sleep(delay)
				if delay < maxDelay {
					delay *= 2
					if delay > maxDelay {
						delay = maxDelay
					}
				}
			} else {
				// reset backoff when connected
				delay = minDelay
			}
		}
	}
}

func newDHT(ctx context.Context, h host.Host, ds datastore.Batching) (*dualdht.DHT, error) {
	dhtOpts := []dualdht.Option{
		dualdht.DHTOption(dht.NamespacedValidator("pk", record.PublicKeyValidator{})),
		dualdht.DHTOption(dht.NamespacedValidator("ipns", ipns.Validator{KeyBook: h.Peerstore()})),
		// dualdht.DHTOption(dht.Concurrency(500)),
		dualdht.DHTOption(dht.Mode(dht.ModeAuto)),
	}
	if ds != nil {
		dhtOpts = append(dhtOpts, dualdht.DHTOption(dht.Datastore(ds)))
	}
	return dualdht.New(ctx, h, dhtOpts...)
}

// GetConnectedPeers returns the list of connected peers
func ConnectedPeers() []*peer.AddrInfo {
	var pinfos []*peer.AddrInfo = []*peer.AddrInfo{}
	host, _ := GetP2PNode(nil)
	for _, p := range host.Peerstore().Peers() {
		// check if the peer is connected
		if host.Network().Connectedness(p) == network.Connected {
			pinfos = append(pinfos, &peer.AddrInfo{
				ID:    p,
				Addrs: host.Peerstore().Addrs(p),
			})
		}
	}
	return pinfos
}

func AllPeers() []*PeerWithStatus {
	var pinfos []*PeerWithStatus = []*PeerWithStatus{}
	host, _ := GetP2PNode(nil)
	for _, p := range host.Peerstore().Peers() {
		pinfos = append(pinfos, &PeerWithStatus{
			ID:            p.String(),
			Connectedness: host.Network().Connectedness(p).String(),
		})
	}
	return pinfos
}

func ConnectedBootstraps() []string {
	var bootstraps = []string{}
	dnt := GetNodeTable()
	host, _ := GetP2PNode(nil)
	for _, p := range *dnt {
		if p.PublicAddress != "" {
			common.Logger.Info("Peer: ", p.ID, " Public Address: ", p.PublicAddress, " Connectedness: ", host.Network().Connectedness(peer.ID(p.ID)), " Host ID: ", host.ID())
			if host.Network().Connectedness(peer.ID(p.ID)) == network.Connected || host.ID().String() == p.ID {
				bootstrapAddr := "/ip4/" + p.PublicAddress + "/tcp/" + viper.GetString("tcpport") + "/p2p/" + p.ID
				bootstraps = append(bootstraps, bootstrapAddr)
			}
		}
	}
	// add myself as bootstrap
	myaddr := host.Addrs()[0].String() + "/p2p/" + host.ID().String()
	bootstraps = append(bootstraps, myaddr)
	// deduplicate
	bootstraps = common.DeduplicateStrings(bootstraps)
	return bootstraps
}

// GetResourceManagerStats returns current resource usage statistics
func GetResourceManagerStats() {
	host, _ := GetP2PNode(nil)
	if rm := host.Network().ResourceManager(); rm != nil {
		// Try to get stats if available
		if statsGetter, ok := rm.(interface {
			Stat() rcmgr.ResourceManagerStat
		}); ok {
			stats := statsGetter.Stat()
			common.Logger.Infof("Resource Manager Stats - System: Conns=%d (in:%d out:%d), Streams=%d (in:%d out:%d), Memory=%d",
				stats.System.NumConnsInbound+stats.System.NumConnsOutbound,
				stats.System.NumConnsInbound,
				stats.System.NumConnsOutbound,
				stats.System.NumStreamsInbound+stats.System.NumStreamsOutbound,
				stats.System.NumStreamsInbound,
				stats.System.NumStreamsOutbound,
				stats.System.Memory,
			)
		} else {
			common.Logger.Info("Resource Manager present but stats not available")
		}
	} else {
		common.Logger.Info("No Resource Manager configured")
	}
}
