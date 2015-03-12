package gae

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine/urlfetch"
	"net/http"
)

func OAuthClient(c context.Context, scope ...string) *http.Client {
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(c, scope...),
			Base:   &urlfetch.Transport{Context: c},
		},
	}
	return hc
}
