// Package interactor contains all use-case interactors preformed by the application.
package interactor

import (
	"context"

	"github.com/lroman242/redirector/domain/dto"
	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/domain/repository"
)

//mockgen -package=mocks -destination=mocks/mock_click_handler.go -source=domain/interactor/click_handler.go ClickHandlerInterface
//go:generate mockgen -package=mocks -destination=mocks/mock_click_handler.go -source=click_handler.go ClickHandlerInterface

// ClickHandlerInterface defines how click events should be processed.
// Implementations can store clicks, forward them to analytics systems,
// or perform other tracking-related operations.
type ClickHandlerInterface interface {
	// HandleClick processes a click event asynchronously and returns a channel
	// that will receive the processing result. This allows for non-blocking
	// click processing while still providing feedback on the operation.
	HandleClick(ctx context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult
}

// ClickHandlerFunc type is a simple implementation of ClickHandlerInterface.
type ClickHandlerFunc func(ctx context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult

// HandleClick function will do some work with the provided entity.Click.
func (ch ClickHandlerFunc) HandleClick(ctx context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult {
	return ch(ctx, click)
}

// storeClickHandler implements ClickHandlerInterface to persist click events to storage
type storeClickHandler struct {
	repo repository.ClicksRepository
}

// NewStoreClickHandler function creates implementation of ClickHandlerInterface
// which saves entity.Click to the storage using repository.ClicksRepository.
func NewStoreClickHandler(clkRepository repository.ClicksRepository) ClickHandlerInterface {
	return &storeClickHandler{repo: clkRepository}
}

// HandleClick function will do some work with the provided entity.Click.
func (sch *storeClickHandler) HandleClick(ctx context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult {
	output := make(chan *dto.ClickProcessingResult)

	go func(ctx context.Context, click *entity.Click) {
		defer close(output)
		// Add context cancellation handling
		select {
		case <-ctx.Done():
			output <- &dto.ClickProcessingResult{
				Click: click,
				Err:   ctx.Err(),
			}
			return
		default:
			output <- &dto.ClickProcessingResult{
				Click: click,
				Err:   sch.repo.Save(ctx, click),
			}
		}
	}(ctx, click)

	return output
}
