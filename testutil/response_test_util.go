package testutil

import (
	"io/ioutil"
	"fmt"
	"regexp"
	"runtime"
	"testing"
	"net/http"
)

var (
	memoizedResponseBodies = make(map[*http.Response]string)
)

func AssertResponseStatus(t *testing.T, expectedStatus int, response *http.Response) {
	if response.StatusCode != expectedStatus {
		fatal(t, "expected %v; got %v. Body: '%v'", expectedStatus, response.StatusCode, ReadResponseBody(t, response))
	}
}

func AssertResponseBody(t *testing.T, expectedBody string, response *http.Response) {
	if body := ReadResponseBody(t, response); body != expectedBody {
		fatal(t, "expected '%v'; got '%v'", expectedBody, body)
	}
}

func AssertResponseBodyMatchesPattern(t *testing.T, expectedPattern string, response *http.Response) {
	if match, body := assertResponseBodyMatch(t, expectedPattern, response); !match {
		fatal(t, "expected body to match '%v'; it did not. Body: '%v'", expectedPattern, body)
	}
}

func AssertResponseBodyDoesNotMatchPattern(t *testing.T, expectedPattern string, response *http.Response) {
	if match, body := assertResponseBodyMatch(t, expectedPattern, response); match {
		fatal(t, "expected body not to match '%v'; it did. Body: '%v'", expectedPattern, body)
	}
}

func assertResponseBodyMatch(t *testing.T, expectedPattern string, response *http.Response) (bool, string) {
	body := ReadResponseBody(t, response)
	if match, err := regexp.MatchString(expectedPattern, body); err != nil {
		fatal(t, "error matching body against pattern %v: %v. Stack: %v", expectedPattern, err)
		return false, ""
	} else {
		return match, body
	}
}

func fatal(t *testing.T, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	t.Fatalf(fmt.Sprintf("%v\n%v", msg, stack()))
}

func stack() string {
	trace := make([]byte, 1024)
	runtime.Stack(trace, false)
	return fmt.Sprintf("%s", trace)
}

func AssertResponseHeaders(t *testing.T, expectedHeaders map[string]string, response *http.Response) {
	if len(expectedHeaders) != len(response.Header) {
		fatal(t, "expected %v headers; got %v. Returned headers: %v", len(expectedHeaders), len(response.Header), response.Header)
	}
	for key, _ := range expectedHeaders {
		if val := response.Header.Get(key); val != expectedHeaders[key] {
			fatal(t, "expected header %v to be '%v'; got '%v'", key, expectedHeaders[key], val)
		}
	}
}

func ReadResponseBody(t *testing.T, response *http.Response) string {
	body, exists := memoizedResponseBodies[response]
	if exists {
		return body
	} else {
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fatal(t, "Error reading body: %v", err)
		}
		body := string(bytes)
		memoizedResponseBodies[response] = body
		return body
	}
}
