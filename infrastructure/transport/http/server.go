package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lroman242/redirector/domain/interactor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewHandler(interactor interactor.RedirectInteractor) http.Handler {
	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/{slug}", NewRedirectHandler(interactor))

	return r
}
