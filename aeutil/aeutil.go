package aeutil

import (
	"net/http"
	"appengine"
	"appengine/aetest"
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

/**
 * The Context returned from this function is registered with the current ContextProvider,
 * so calls to UncurryWithContext will resolve to this context.
 */
func CreateAndRegisterTestContext() (context aetest.Context, err error) {
	context, err = aetest.NewContext(nil)

	// TODO change this function to maintain and consult a map, so we can
	// make multiple requests in parallel / in a recursive fashion
	ContextProvider = func(*http.Request) appengine.Context { return context }

	return
}
