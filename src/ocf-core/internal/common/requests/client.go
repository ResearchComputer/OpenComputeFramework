package requests

import (
	"github.com/sethgrid/pester"
)

var client *pester.Client

func NewHTTPClient() *pester.Client {
	if client == nil {
		client = pester.New()
		client.MaxRetries = 1
		client.Concurrency = 1
		client.Backoff = pester.ExponentialJitterBackoff
	}
	return client
}
