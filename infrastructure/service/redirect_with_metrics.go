// Package service provides implementations of domain service interfaces and additional
// service wrappers for metrics, caching, and other cross-cutting concerns.
package service

import (
	"context"
	"time"

	"github.com/lroman242/redirector/domain/dto"
	"github.com/lroman242/redirector/domain/interactor"
	"github.com/lroman242/redirector/infrastructure/metrics"
)

// RedirectWithMetrics is a wrapper around the RedirectInteractor that adds metrics tracking.
// It implements the RedirectInteractor interface and decorates the redirect operations
// with Prometheus metrics for monitoring and observability.
type RedirectWithMetrics struct {
	// RedirectInteractor is the wrapped redirect interactor implementation
	RedirectInteractor interactor.RedirectInteractor
}

// NewRedirectWithMetrics creates a new RedirectWithMetrics instance that wraps
// the provided RedirectInteractor with metrics tracking functionality.
func NewRedirectWithMetrics(redirectInteractor interactor.RedirectInteractor) *RedirectWithMetrics {
	return &RedirectWithMetrics{
		RedirectInteractor: redirectInteractor,
	}
}

// Redirect handles redirect requests and tracks metrics about the operation.
// It increments counters for total redirects and redirects by slug, and measures
// the execution time of the redirect operation.
func (r *RedirectWithMetrics) Redirect(ctx context.Context, slug string, requestData *dto.RedirectRequestData) (*dto.RedirectResult, error) {
	// Track total redirects and redirects by slug
	metrics.RedirectTotal.Inc()
	metrics.RedirectsBySlug.WithLabelValues(slug).Inc()

	// Measure execution time
	startTime := time.Now()
	defer func() {
		metrics.RedirectDuration.Observe(time.Since(startTime).Seconds())
	}()

	return r.RedirectInteractor.Redirect(ctx, slug, requestData)
}
