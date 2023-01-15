package daemon

import (
	"rccore/internal/common"
	"rccore/internal/server"

	"github.com/getsentry/sentry-go"
)

func Start() {
	common.Logger.Info("DSN: " + common.BuildSecret.SentryDSN)
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

	StartTicker()
	server.StartServer()
}
