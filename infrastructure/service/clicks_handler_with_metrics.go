// Package service provides implementations of domain service interfaces and additional
// service wrappers for metrics, caching, and other cross-cutting concerns.
package service

import (
	"context"
	"time"

	"github.com/lroman242/redirector/domain/dto"
	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/domain/interactor"
	"github.com/lroman242/redirector/infrastructure/metrics"
)

// ClickHandlerWithMetrics wraps a ClickHandlerInterface implementation
// and adds Prometheus metrics tracking for click processing duration.
type ClickHandlerWithMetrics struct {
	// Handler is the wrapped click handler implementation
	interactor.ClickHandlerInterface
	// handlerName is used as a label in metrics to identify the handler
	handlerName string
}

// NewClickHandlerWithMetrics creates a new ClickHandlerWithMetrics instance
// that decorates the provided handler with metrics tracking functionality.
func NewClickHandlerWithMetrics(handler interactor.ClickHandlerInterface, handlerName string) interactor.ClickHandlerInterface {
	return &ClickHandlerWithMetrics{handler, handlerName}
}

// HandleClick processes a click event and tracks its duration using Prometheus metrics.
// It delegates the actual click processing to the wrapped handler while measuring
// the time taken to process the click.
func (h *ClickHandlerWithMetrics) HandleClick(ctx context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult {
	start := time.Now()
	defer func() {
		metrics.ClickHandlerDuration.WithLabelValues(h.handlerName).Observe(time.Since(start).Seconds())
	}()

	return h.ClickHandlerInterface.HandleClick(ctx, click)
}
