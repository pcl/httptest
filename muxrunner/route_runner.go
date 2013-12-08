package muxrunner

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"bufio"
	"bytes"
	"reflect"
	mockHttp "github.com/stretchr/testify/http"
)

type RequestRunner interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewInProcessRequestRunner() RequestRunner {
	return NewInProcessRequestRunnerWithServeMux(http.DefaultServeMux)
}

func NewInProcessRequestRunnerWithServeMux(mux *http.ServeMux) RequestRunner {
	return InProcessRequestRunner{ mux: mux }
}

type InProcessRequestRunner struct {
	mux *http.ServeMux
}

func (r InProcessRequestRunner) Do(req *http.Request) (response *http.Response, err error) {
	var match mux.RouteMatch
	switch handler, _ := r.mux.Handler(req); handler.(type) {
	default:
		err = fmt.Errorf("handler for request '%v' was not of type %v", req.RequestURI, reflect.TypeOf((*mux.Router)(nil)))
		return

	case *mux.Router:
		if !handler.(*mux.Router).Match(req, &match) {
			err = fmt.Errorf("no match found for path '%v'", req.RequestURI)
		}
	}

	responseWriter := mockHttp.TestResponseWriter{}
	match.Handler.(http.HandlerFunc)(&responseWriter, req) // TODO check that it's a HandlerFunc to avoid a panic

	var httpBuffer bytes.Buffer
	fmt.Fprintf(&httpBuffer, "HTTP/1.1 %v\r\n", statusLine(responseWriter.WrittenHeaderInt))
	responseWriter.Header().Write(&httpBuffer)
	fmt.Fprintf(&httpBuffer, "\r\n%v", responseWriter.Output)

	response, err = http.ReadResponse(bufio.NewReader(&httpBuffer), req)
	return
}

func statusLine(status int) string {
	if status == 0 {
		status = 200
	}
	return fmt.Sprintf("%v", status) // TODO do we need to add the human-readable form to be compliant?
}
