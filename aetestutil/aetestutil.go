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
	aeutil.ContextProvider = func(req *http.Request) appengine.Context {
		// Copy the request details from the user-created one to the
		// context's one, and edit it to add in the Appengine headers,
		// as seen in appengine/user/user_dev.go. This has to happen
		// before we make any calls to the user module, which works
		// out since this is called during UncurryWithContext.
		contextRequest := context.Request().(*http.Request)
		contextRequest.Body = req.Body
		contextRequest.ContentLength = req.ContentLength
		contextRequest.Form = req.Form
		contextRequest.Header = req.Header
		contextRequest.Host = req.Host
		contextRequest.Method = req.Method
		contextRequest.MultipartForm = req.MultipartForm
		contextRequest.PostForm = req.PostForm
		contextRequest.Proto = req.Proto
		contextRequest.ProtoMajor = req.ProtoMajor
		contextRequest.ProtoMinor = req.ProtoMinor
		contextRequest.RemoteAddr = req.RemoteAddr
		contextRequest.RequestURI = req.RequestURI
		contextRequest.TLS = req.TLS
		contextRequest.Trailer = req.Trailer
		contextRequest.TransferEncoding = req.TransferEncoding
		contextRequest.URL = req.URL

		return context
	}

	return
}
