package daemon

import (
	"ocfcore/internal/common"
	"ocfcore/internal/server"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
)

func StartTicker() {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(viper.GetInt("vacuum.interval")).Seconds().Do(func() {
		common.Logger.Debug("Vacuuming...")
		server.DisconnectionDetection(time.Duration(viper.GetInt("vacuum.tolerance")) * time.Second)

	})
	if err != nil {
		common.Logger.Error(err)
	}
	s.StartAsync()
}
