// Package click_handlers provides implementations of click event handlers
// with additional functionality like metrics tracking.
package click_handlers

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
	// ClickHandlerInterface is the wrapped click handler implementation.
	interactor.ClickHandlerInterface
}

// NewClickHandlerWithMetrics creates a new ClickHandlerWithMetrics instance
// that decorates the provided handler with metrics tracking functionality.
func NewClickHandlerWithMetrics(handler interactor.ClickHandlerInterface) interactor.ClickHandlerInterface {
	return &ClickHandlerWithMetrics{handler}
}

// HandleClick processes a click event and tracks its duration using Prometheus metrics.
// It delegates the actual click processing to the wrapped handler while measuring
// the time taken to process the click.
func (h *ClickHandlerWithMetrics) HandleClick(ctx context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult {
	handlerName := "store_click_handler"
	start := time.Now()
	defer func() {
		metrics.ClickHandlerDuration.WithLabelValues(handlerName).Observe(time.Since(start).Seconds())
	}()

	return h.ClickHandlerInterface.HandleClick(ctx, click)
}
