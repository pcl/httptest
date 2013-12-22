package aetestutil

import (
	"net/http"
	"appengine"
	"appengine/aetest"
	"github.com/pcl/httptest/aeutil"
)

/**
 * The Context returned from this function is registered with the current ContextProvider,
 * so calls to aeutil.UncurryWithContext will resolve to this context.
 */
func CreateAndRegisterTestContext() (context aetest.Context, err error) {
	context, err = aetest.NewContext(nil)

	// TODO change this function to maintain and consult a map, so we can
	// make multiple requests in parallel / in a recursive fashion
	aeutil.ContextProvider = func(*http.Request) appengine.Context { return context }

	return
}
