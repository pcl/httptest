package aetestutil

import (
	"appengine"
	"appengine/aetest"
	"github.com/pcl/httptest/muxrunner"
	"github.com/pcl/httptest/aeutil"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
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

		// ##### should have some way to register what header values to use here
		req.Header.Set("X-AppEngine-Internal-User-Email", "foo@example.com")
		req.Header.Set("X-AppEngine-Internal-User-Id", "foo@example.com")
		return context
	}

	// ##### should have some way to register what header values to use here
	req := context.Request().(*http.Request)
	req.Header.Set("X-AppEngine-Internal-User-Email", "foo@example.com")
	req.Header.Set("X-AppEngine-Internal-User-Id", "foo@example.com")

	return
}

func CreateRegisterAndAssertTestContext(t *testing.T) aetest.Context {
	c, err := CreateAndRegisterTestContext()
	if err != nil {
		t.Fatalf("CreateRegisterAndAssertTestContext: %v", err)
	}
	return c
}

func ExecuteAndAssertPostWithFormValues(t *testing.T, path string, formValues map[string]string) *http.Response {

	req := CreateAndAssertHttpRequest(t, "POST", path)

	values := url.Values{}
	for k, v := range formValues {
		values.Add(k, v)
	}
	req.Form = values

	return doRequestAndAssert(t, req)
}

func ExecuteAndAssertPostWithBodyText(t *testing.T, path string, body string) *http.Response {
	req := CreateAndAssertHttpRequest(t, "POST", path)
	req.Body = ioutil.NopCloser(strings.NewReader(body))
	return doRequestAndAssert(t, req)
}

func ExecuteAndAssertGet(t *testing.T, path string) *http.Response {
	req := CreateAndAssertHttpRequest(t, "GET", path)
	return doRequestAndAssert(t, req)
}

func doRequestAndAssert(t *testing.T, req *http.Request) *http.Response {
	var client *http.Client
	if testHost() == "" {
		client = muxrunner.InProcessClient()
		req.URL.Scheme = "http"
		req.URL.Host = "localhost"
	}

    response, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error executing request against '%v': %v", req.URL, err)
	}
	return response
}

func CreateAndAssertHttpRequest(t *testing.T, method, path string) (req *http.Request) {
	location := testHost() + path
	var err error
	req, err = http.NewRequest(method, location, nil)
	if err != nil {
		t.Fatalf("Error creating request for '%v': %v", location, err)
	}

	return
}

func testHost() string {
	return os.Getenv("TARGET_HOST")
}
