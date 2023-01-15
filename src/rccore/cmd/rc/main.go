package main

import (
	"rccore/cmd/rc/cmd"
	"rccore/internal/common"
)

var (
	// Populated during build
	version      = "dev"
	commitHash   = "?"
	buildDate    = ""
	authUrl      = ""
	authClientId = ""
	authSecret   = ""
	sentryDSN    = ""
)

func main() {
	common.JSONVersion.Version = version
	common.JSONVersion.Commit = commitHash
	common.JSONVersion.Date = buildDate
	common.BuildSecret.AuthClientID = authClientId
	common.BuildSecret.AuthURL = authUrl
	common.BuildSecret.AuthSecret = authSecret
	common.BuildSecret.SentryDSN = sentryDSN
	cmd.Execute()
}
