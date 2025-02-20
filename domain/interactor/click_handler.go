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

// ClickHandlerInterface describes handler which can manage the entity.Click.
type ClickHandlerInterface interface {
	HandleClick(ctx context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult
}

// ClickHandlerFunc type is a simple implementation of ClickHandlerInterface.
type ClickHandlerFunc func(ctx context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult

// HandleClick function will do some work with the provided entity.Click.
func (ch ClickHandlerFunc) HandleClick(ctx context.Context, click *entity.Click) <-chan *dto.ClickProcessingResult {
	return ch(ctx, click)
}

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
		output <- &dto.ClickProcessingResult{
			Click: click,
			Err:   sch.repo.Save(ctx, click),
		}
	}(ctx, click)

	return output
}
