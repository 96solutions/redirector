// Package metrics provides functionality for tracking and exposing application metrics.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Cache metrics
	// TODO:
	CacheOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redirective_cache_operations_total",
			Help: "Total number of cache operations",
		},
		[]string{"slug", "result"}, // slug: incoming slug, result: hit/miss
	)

	// RedirectTotal tracks the total number of handled redirects.
	RedirectTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redirector_redirects_total",
		Help: "The total number of handled redirects.",
	})

	// RedirectsBySlug tracks the number of redirects per slug.
	RedirectsBySlug = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "redirector_redirects_by_slug_total",
		Help: "The total number of redirects handled per slug.",
	}, []string{"slug"})

	// RedirectDuration tracks the execution time of redirect requests.
	RedirectDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "redirector_redirect_duration_seconds",
		Help:    "The time taken to process redirect requests.",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	})

	// ClickHandlerDuration tracks the processing time per click handler.
	ClickHandlerDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "redirector_click_handler_duration_seconds",
		Help:    "The time taken to process clicks by handler.",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"handler"})
)
