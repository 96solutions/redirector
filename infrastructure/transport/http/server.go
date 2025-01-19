package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lroman242/redirector/domain/interactor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewHandler(interactor interactor.RedirectInteractor) http.Handler {
	r := mux.NewRouter()

	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Redirector!"))
	}))
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/r/{slug}", NewRedirectHandler(interactor))

	return r
}
