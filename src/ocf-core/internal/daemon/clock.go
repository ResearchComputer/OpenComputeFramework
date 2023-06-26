package daemon

import (
	"ocfcore/internal/common"
	"ocfcore/internal/server"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
)

var firstRun = true

func StartTicker() {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(viper.GetInt("vacuum.interval")).Seconds().Do(func() {
		if firstRun {
			// skip the first run to wait until the server is ready
			firstRun = false
			return
		}
		common.Logger.Debug("Vacuuming...")
		server.DisconnectionDetection(time.Duration(viper.GetInt("vacuum.tolerance")) * time.Second)
		server.UpdateGlobalWorkloadTable()
	})
	if err != nil {
		common.Logger.Error(err)
	}
	s.StartAsync()
}
