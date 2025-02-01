package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lroman242/redirector/domain/interactor"
)

// NewHandler register a new HTTP handler (router).
func NewHandler(interactor interactor.RedirectInteractor) http.Handler {
	r := mux.NewRouter()

	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Redirector"))
	}))

	r.Handle("/r/{slug}", NewRedirectHandler(interactor))

	return r
}
