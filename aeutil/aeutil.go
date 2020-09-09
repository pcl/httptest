package aeutil

import (
	"context"
	"net/http"
	"google.golang.org/appengine"
)

var (
	// To be set by an external environment, if it wants control over context creation
	ContextProvider func(*http.Request) context.Context
)

func init() {
	ContextProvider = func(r *http.Request) context.Context {
		return appengine.NewContext(r)
	}
}

func UncurryWithContext(handler func(http.ResponseWriter, *http.Request, context.Context)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, ContextProvider(r))
	}
}
