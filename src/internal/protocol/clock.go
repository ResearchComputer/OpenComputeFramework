package protocol

import (
	"context"
	"math/rand"
	"ocf/internal/common"
	"time"

	"github.com/go-co-op/gocron"
	ds "github.com/ipfs/go-datastore"
)

var verificationKey = "ocf-verification-key"
var verificationProb = 0.5

func StartTicker() {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Minute().Do(func() {
		if rand.Float64() < verificationProb {
			store, _ := GetCRDTStore()
			ctx := context.Background()
			store.Put(ctx, ds.NewKey(verificationKey), []byte("verification"))
		}
	})
	common.ReportError(err, "Error while starting ticker")
}
