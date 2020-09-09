package muxrunner

import (
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"fmt"
	"bufio"
	"bytes"
	"reflect"
)

func InProcessClient() *http.Client {
	return InProcessClientWithServeMux(http.DefaultServeMux)
}

func InProcessClientWithServeMux(mux *http.ServeMux) *http.Client {
	return &http.Client{ Transport: inProcessRoundTripper{ mux: mux } }
}

type inProcessRoundTripper struct {
	mux *http.ServeMux
}

func (r inProcessRoundTripper) RoundTrip(req *http.Request) (response *http.Response, err error) {
	var match mux.RouteMatch
	switch handler, _ := r.mux.Handler(req); handler.(type) {
	default:
		err = fmt.Errorf("handler for request '%v' was not of type %v", req.RequestURI, reflect.TypeOf((*mux.Router)(nil)))
		return

	case *mux.Router:
		muxRouter := handler.(*mux.Router)
		if !muxRouter.Match(req, &match) {
			err = fmt.Errorf("no match found for path '%v'", req.URL)
			return
		}

		responseWriter := httptest.NewRecorder()
		muxRouter.ServeHTTP(responseWriter, req)

		var httpBuffer bytes.Buffer
		fmt.Fprintf(&httpBuffer, "HTTP/1.1 %v\r\n", statusLine(responseWriter.Code))
		responseWriter.Header().Write(&httpBuffer)
		fmt.Fprintf(&httpBuffer, "\r\n%v", responseWriter.Body.String())

		response, err = http.ReadResponse(bufio.NewReader(&httpBuffer), req)
	}

	return
}

func statusLine(status int) string {
	if status == 0 {
		status = 200
	}
	return fmt.Sprintf("%v", status) // TODO do we need to add the human-readable form to be compliant?
}
