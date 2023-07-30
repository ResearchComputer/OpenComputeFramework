package daemon

import (
	"ocfcore/internal/common"
	"ocfcore/internal/server"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
)

func Start() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: common.BuildSecret.SentryDSN,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		common.Logger.Error("sentry.Init: %s", err)
	}
	common.Logger.Info("Wallet: ", viper.Get("wallet.account"))
	StartTicker()
	server.StartServer()
}
