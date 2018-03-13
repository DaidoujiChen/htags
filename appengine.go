package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var isInAppEngine bool

func init() {
	isInAppEngine = true
}

func startService(handler func(http.ResponseWriter, *http.Request)) {
	if isInAppEngine {
		fmt.Println("Start in App Engine")
		http.HandleFunc("/", indexHandler)
		appengine.Main()
	} else {
		fmt.Println("Start for Develop")
		http.HandleFunc("/", indexHandler)
		http.ListenAndServe(":8030", nil)
	}
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
