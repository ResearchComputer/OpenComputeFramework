package main

import (
	"ocf/bin/cmd"
	"ocf/internal/common"
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
	cmd.Execute()
}
