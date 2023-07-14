package remote

import (
	"io"
	"strings"

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

func HTTPPost(remoteAddr string, req []byte) (string, error) {
	payload := strings.NewReader(string(req))
	resp, err := NewHTTPClient().Post(remoteAddr, "application/json", payload)
	if err != nil {
		return "nil", err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(b), err
	}
	return string(b), nil
}
