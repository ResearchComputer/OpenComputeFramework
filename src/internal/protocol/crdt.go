package protocol

import (
	"context"
	"fmt"
	"ocf/internal/common"
	"sync"
	"time"

	ipfslite "github.com/hsanjuan/ipfs-lite"
	ds "github.com/ipfs/go-datastore"
	badger "github.com/ipfs/go-ds-badger"
	crdt "github.com/ipfs/go-ds-crdt"
	gocrdt "github.com/ipfs/go-ds-crdt"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	pubsubTopic = "ocf-crdt"
	pubsubKey   = "ocf-crdt"
)
var crdtStore *gocrdt.Datastore
var once sync.Once
var cancelSubscriptions context.CancelFunc

func GetCRDTStore() (*gocrdt.Datastore, context.CancelFunc) {
	once.Do(func() {
		host, dht := GetP2PNode(nil)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		store, err := badger.NewDatastore(common.GetDBPath(), &badger.DefaultOptions)
		common.ReportError(err, "Error while creating datastore")

		ipfs, err := ipfslite.New(ctx, store, nil, host, &dht, nil)
		common.ReportError(err, "Error while creating ipfs lite node")

		psub, err := pubsub.NewGossipSub(ctx, host)
		common.ReportError(err, "Error while creating pubsub")

		topic, err := psub.Join(pubsubTopic)
		common.ReportError(err, "Error while joining pubsub topic")

		netSubs, err := topic.Subscribe()
		common.ReportError(err, "Error while subscribing to pubsub topic")

		go func() {
			for {
				msg, err := netSubs.Next(ctx)
				if err != nil {
					fmt.Println(err)
					break
				}
				host.ConnManager().TagPeer(msg.ReceivedFrom, "keep", 100)
			}
		}()

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					topic.Publish(ctx, []byte("ping"))
					time.Sleep(20 * time.Second)
				}
			}
		}()
		psubCtx, pcancel := context.WithCancel(ctx)
		cancelSubscriptions = pcancel
		pubsubBC, err := crdt.NewPubSubBroadcaster(psubCtx, psub, pubsubTopic)
		common.ReportError(err, "Error while creating pubsub broadcaster")
		opts := crdt.DefaultOptions()
		opts.Logger = common.Logger
		opts.RebroadcastInterval = 5 * time.Second
		opts.PutHook = func(k ds.Key, v []byte) {
			fmt.Printf("Added: [%s] -> %s\n", k, string(v))
		}
		opts.DeleteHook = func(k ds.Key) {
			fmt.Printf("Removed: [%s]\n", k)
		}

		crdtStore, err = crdt.New(store, ds.NewKey(pubsubKey), ipfs, pubsubBC, opts)
		common.ReportError(err, "Error while creating crdt store")

		common.Logger.Info("Bootstrapping...")
		addsInfo, err := peer.AddrInfosFromP2pAddrs(getDefaultBootstrapPeers(nil)...)
		common.ReportError(err, "Error while getting bootstrap peers")
		ipfs.Bootstrap(addsInfo)
		// h.ConnManager().TagPeer(inf.ID, "keep", 100)
	})
	return crdtStore, cancelSubscriptions
}
