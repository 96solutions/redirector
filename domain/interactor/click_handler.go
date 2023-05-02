package interactor

import (
	"context"
	"github.com/lroman242/redirector/domain/entity"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_click_handler.go -source=domain/interactor/click_handler.go ClickHandlerInterface
// ClickHandlerInterface describes handler which can manage the entity.Click.
type ClickHandlerInterface interface {
	HandleClick(ctx context.Context, click *entity.Click) <-chan *clickProcessingResult
}

// ClickHandlerFunc type is a simple implementation of ClickHandlerInterface.
type ClickHandlerFunc func(ctx context.Context, click *entity.Click) <-chan *clickProcessingResult

// HandleClick function will do some work with the provided entity.Click.
func (ch ClickHandlerFunc) HandleClick(ctx context.Context, click *entity.Click) <-chan *clickProcessingResult {
	return ch(ctx, click)
}
