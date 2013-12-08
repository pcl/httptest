## HTTP Test Utilities ##

This module contains a number of Go HTTP test utilities.

1. [In-process gorilla/mux test harness](#in-process-gorillamux-test-harness)
2. [Unit testing utilities (assertions etc.)](#http-testing-utilities)

### In-process gorilla/mux test harness ###

The [Gorilla mux](http://www.gorillatoolkit.org/pkg/mux) package is a powerful framework
for building RESTful HTTP services. However, it's not as testable as it could be. In particular,
tests written for it must either be written directly against the handler methods themselves, or
an HTTP server must be externally brought online and configured for the test environment. The
muxrunner package provides some functions to write endpoint-oriented tests that can run in a
 variety of different contexts, including in-process testing during development and remote
 testing for pre-deploy or ongoing system verification.

(Currently, only the in-process variant is implemented. Next step: implement the remote testing
 function, and environment variables to detect the appropriate current context.)

#### Usage ####

First, set up your handlers as you would normally:

```go
func init() {
    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    http.Handle("/", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "HomeHandler")
}
```

Next, write a test that creates a RequestRunner and executes a request against it:

```go
func TestHomeHandler() {
    runner := muxrunner.NewInProcessRequestRunner()
    req, _ := http.NewRequest("GET", "/path/to/handler", nil)
    response, _ := runner.Do(req)
    testutil.AssertResponseStatus(t, 200, response)
    testutil.AssertResponseBody(t, "HomeHandler", response)
}
```

### HTTP testing utilities ###

This package provides some useful functions for writing HTTP-oriented tests,
including response body and header assertions. See the [usage for muxrunner](#usage)
above for an example.