package test

import (
	"net/http"
	"testing"
	"github.com/pcl/httptest/muxrunner"
	"github.com/pcl/httptest/testutil"
)

func TestHomeHandler(t *testing.T) {
	runner := muxrunner.NewInProcessRequestRunner()
	req, _ := http.NewRequest("GET", "/", nil)
	response, _ := runner.Do(req)
	testutil.AssertResponseStatus(t, 200, response)
	testutil.AssertResponseBody(t, "HomeHandler", response)
}

func TestHeaders(t *testing.T) {
	runner := muxrunner.NewInProcessRequestRunner()
	req, _ := http.NewRequest("GET", "/headers", nil)
	req.Header.Add("Foo", "Bar")
	response, _ := runner.Do(req)
	testutil.AssertResponseStatus(t, 200, response)
	testutil.AssertResponseBody(t, "HeaderHandler", response)
	testutil.AssertResponseHeaders(t, map[string]string { "Foo": "Bar" }, response)
}
