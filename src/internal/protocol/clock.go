package protocol

import (
	"math/rand"
	"ocf/internal/common"

	"github.com/jasonlvhit/gocron"
)

// var verificationKey = "ocf-verification-key"
var verificationProb = 0.5

func StartTicker() {
	err := gocron.Every(1).Minute().Do(func() {
		common.Logger.Info("Starting verification")
		if rand.Float64() < verificationProb {
			// store, _ := GetCRDTStore()
			// ctx := context.Background()
			// store.Put(ctx, ds.NewKey(verificationKey), []byte("verification"))
			Reconnect()
		}
	})
	<-gocron.Start()
	common.ReportError(err, "Error while starting ticker")
}
