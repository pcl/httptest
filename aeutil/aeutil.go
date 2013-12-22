package aeutil

import (
	"net/http"
	"appengine"
)

var (
	// To be set by an external environment, if it wants control over context creation
	ContextProvider func(*http.Request) appengine.Context
)

func init() {
	ContextProvider = func(r *http.Request) appengine.Context {
		return appengine.NewContext(r)
	}
}

func UncurryWithContext(handler func(http.ResponseWriter, *http.Request, appengine.Context)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, ContextProvider(r))
	}
}
