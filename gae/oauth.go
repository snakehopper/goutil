package gae

import (
	"appengine"
	"github.com/golang/oauth2/google"
	"net/http"
)

func OAuthClient(c appengine.Context, scope []string) *http.Client {
	config := google.NewAppEngineConfig(c, scope)

	return &http.Client{Transport: config.NewTransport()}
}
