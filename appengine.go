package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var isInAppEngine bool

func init() {
	isInAppEngine = true
}

func httpGet(url string, r *http.Request) (*http.Response, error) {
	if isInAppEngine {
		ctx := appengine.NewContext(r)
		client := urlfetch.Client(ctx)
		return client.Get(url)
	} else {
		return http.Get(url)
	}
}
