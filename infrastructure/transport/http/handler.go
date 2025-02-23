// Package http provides HTTP transport layer implementations.
package http

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lroman242/redirector/domain/dto"
	"github.com/lroman242/redirector/domain/interactor"
	"github.com/lroman242/redirector/infrastructure/metrics"
	uuid "github.com/satori/go.uuid"
)

// RedirectHandler handles HTTP redirect requests by delegating to a RedirectInteractor.
type RedirectHandler struct {
	interactor interactor.RedirectInteractor
}

// NewRedirectHandler creates a new RedirectHandler instance.
func NewRedirectHandler(interactor interactor.RedirectInteractor) *RedirectHandler {
	return &RedirectHandler{interactor: interactor}
}

// getIPAddress extracts the real client IP address from request headers.
// It checks various headers in order of reliability:
// 1. CF-Connecting-IP (Cloudflare)
// 2. X-Real-IP
// 3. X-Forwarded-For (first IP in the chain)
// 4. RemoteAddr as fallback
func getIPAddress(r *http.Request) (net.IP, error) {
	// Check Cloudflare header first
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		if ip := net.ParseIP(cfIP); ip != nil {
			return ip, nil
		}
	}

	// Check X-Real-IP header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		if ip := net.ParseIP(realIP); ip != nil {
			return ip, nil
		}
	}

	// Check X-Forwarded-For header
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			if ip := net.ParseIP(strings.TrimSpace(ips[0])); ip != nil {
				return ip, nil
			}
		}
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	return userIP, nil
}

// ServeHTTP handles HTTP redirect requests.
// It extracts request parameters, calls the redirect interactor,
// and performs the redirect while tracking metrics.
func (rh *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.RedirectDuration.Observe(float64(time.Since(start).Milliseconds()))
	}()

	// Extract slug from URL path
	vars := mux.Vars(r)
	slug := vars["slug"]
	if slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	// Get real client IP address
	userIP, err := getIPAddress(r)
	if err != nil {
		slog.Error("Failed to get client IP", slog.String("error", err.Error()))
		http.Error(w, "Failed to process client IP", http.StatusInternalServerError)
		return
	}

	// Prepare redirect request data
	data := &dto.RedirectRequestData{
		Slug:      slug,
		Params:    r.URL.Query(),
		Headers:   r.Header,
		UserAgent: r.UserAgent(),
		IP:        userIP,
		Protocol:  r.Proto,
		Referer:   r.Referer(),
		URL:       r.URL,
		RequestID: uuid.NewV4().String(),
	}

	slog.Debug("Redirect request", slog.String("slug", slug), "data", data)

	// Process redirect
	redirectResult, err := rh.interactor.Redirect(r.Context(), slug, data)
	if err != nil {
		slog.Error("Redirect failed", slog.String("error", err.Error()), slog.String("slug", slug), slog.Any("request", data))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Debug("redirect result", slog.String("slug", slug), "data", data, "redirectResult", redirectResult)

	// Track metrics
	metrics.RedirectTotal.Inc()
	metrics.RedirectsBySlug.WithLabelValues(slug).Inc()

	// Perform redirect
	http.Redirect(w, r, redirectResult.TargetURL, http.StatusSeeOther)

	// Process click results asynchronously
	go func() {
		for {
			select {
			case <-r.Context().Done():
				return
			case result, ok := <-redirectResult.OutputCh:
				slog.Debug("redirect result", slog.String("slug", slug), "result", result, slog.Bool("isClosed", !ok))
				if !ok {
					slog.Debug("Click processing complete", slog.String("slug", slug))
					return
				}
				if result.Err != nil {
					slog.Error("Click processing failed",
						slog.String("error", result.Err.Error()),
						slog.String("slug", slug),
						slog.Any("request", data),
					)
				}
			}
		}
	}()
}
