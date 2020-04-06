package httplog_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crhntr/httplog"
)

func ExampleWrap() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, world!")
	})

	logMux := httplog.Wrap(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	logMux.ServeHTTP(w, r)
}
// another pull request
func TestWrap(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, world!")
	})

	logMux := httplog.Wrap(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	logMux.ServeHTTP(w, r)
}
