package common

type buildSecret struct {
	AuthClientID string `json:"AUTH_CLIENT_ID"`
	AuthURL      string `json:"AUTH_URL"`
	AuthSecret   string `json:"AUTH_CLIENT_SECRET"`
	SentryDSN    string `json:"SENTRY_DSN"`
}

var BuildSecret buildSecret
