package testutil

import (
	"io/ioutil"
	"testing"
	"net/http"
)

func AssertResponseStatus(t *testing.T, expectedStatus int, response *http.Response) {
	if response.StatusCode != expectedStatus {
		t.Fatalf("expected 200; got '%v' Body: '%v'", response.StatusCode, ReadResponseBody(t, response))
	}
}

func AssertResponseBody(t *testing.T, expectedBody string, response *http.Response) {
	if body := ReadResponseBody(t, response); body != expectedBody {
		t.Fatalf("expected '%v'; got '%v'", expectedBody, body)
	}
}

func AssertResponseHeaders(t *testing.T, expectedHeaders map[string]string, response *http.Response) {
	if len(expectedHeaders) != len(response.Header) {
		t.Fatalf("expected %v headers; got %v. Returned headers: %v", len(expectedHeaders), len(response.Header), response.Header)
	}
	for key, _ := range expectedHeaders {
		if val := response.Header.Get(key); val != expectedHeaders[key] {
			t.Fatalf("expected header %v to be '%v'; got '%v'", key, expectedHeaders[key], val)
		}
	}
}

func ReadResponseBody(t *testing.T, response *http.Response) string {
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Error reading body: %v", err)
	}
	return string(bytes)
}
