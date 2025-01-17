package http

import (
	"github.com/gorilla/mux"
	"github.com/lroman242/redirector/domain/dto"
	"github.com/lroman242/redirector/domain/interactor"
	"log/slog"
	"net"
	"net/http"
)

type RedirectHandler struct {
	interactor interactor.RedirectInteractor
}

func NewRedirectHandler(interactor interactor.RedirectInteractor) *RedirectHandler {
	return &RedirectHandler{interactor: interactor}
}

func (rh *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO: handle Cloudflare, Proxy, etc
	vars := mux.Vars(r)

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		slog.Error(err.Error(), "ip", ip)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		slog.Error("invalid IP address provided", "ip", ip)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := &dto.RedirectRequestData{
		Params:    r.URL.Query(),
		Headers:   r.Header,
		UserAgent: r.UserAgent(),
		IP:        userIP,
		Protocol:  r.Proto,
		Referer:   r.Referer(),
		//RequestID:
	}

	slug := vars["slug"]

	redirectResult, err := rh.interactor.Redirect(r.Context(), slug, data)
	if err != nil {
		slog.Error(err.Error(), "slug", slug, "request", data)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectResult.TargetURL, http.StatusSeeOther)

	go func() {
		for {
			select {
			case <-r.Context().Done():
				break
			case result, cl := <-redirectResult.OutputCh:
				if cl {
					break
				}

				if result.Err != nil {
					slog.Error(result.Err.Error(), "slug", slug, "request", data)
				}
				break
			default:
			}
		}
	}()

	return
}
