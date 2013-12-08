package test

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

// standard web handler registration
func init() {
	r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/headers", HeaderHandler)
	http.Handle("/", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HomeHandler")
}

func HeaderHandler(w http.ResponseWriter, r *http.Request) {
	for key, values := range r.Header {
		for i := 0; i < len(values); i++ {
			w.Header().Add(key, values[i])
		}
	}
	fmt.Fprintf(w, "HeaderHandler")
}
