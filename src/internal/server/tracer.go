package server

import (
	"context"
	"log"
	"ocf/internal/common"
	"os"
	"sync"

	"github.com/axiomhq/axiom-go/axiom"
	axiotel "github.com/axiomhq/axiom-go/axiom/otel"
)

var (
	dataset    = os.Getenv("AXIOM_DATASET")
	tracerOnce sync.Once
	tracker    *axiom.Client
)

func initTracer() {
	tracerOnce.Do(func() {
		ctx := context.Background()
		stop, err := axiotel.InitTracing(ctx, dataset, "research-computer-coordinator", "v1.0.0")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if stopErr := stop(); stopErr != nil {
				log.Fatal(stopErr)
			}
		}()
		tracker, err = axiom.NewClient()
		if err != nil {
			log.Fatal(err)
		}
	})
}

func IngestEvents(events []axiom.Event) {

	go func() {
		// expand events to axiom.Event
		res, err := tracker.IngestEvents(context.Background(), dataset, events)
		if err != nil {
			common.Logger.Error(err)
		}
		for _, fail := range res.Failures {
			common.Logger.Error(fail)
		}
	}()
}
