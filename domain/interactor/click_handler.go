package interactor

import (
	"context"

	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/domain/repository"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_click_handler.go -source=domain/interactor/click_handler.go ClickHandlerInterface
// ClickHandlerInterface describes handler which can manage the entity.Click.
type ClickHandlerInterface interface {
	HandleClick(ctx context.Context, click *entity.Click) <-chan *ClickProcessingResult
}

// ClickHandlerFunc type is a simple implementation of ClickHandlerInterface.
type ClickHandlerFunc func(ctx context.Context, click *entity.Click) <-chan *ClickProcessingResult

// HandleClick function will do some work with the provided entity.Click.
func (ch ClickHandlerFunc) HandleClick(ctx context.Context, click *entity.Click) <-chan *ClickProcessingResult {
	return ch(ctx, click)
}

type storeClickHandler struct {
	repo repository.ClicksRepository
}

func NewStoreClickHandler(clkRepository repository.ClicksRepository) ClickHandlerInterface {
	return &storeClickHandler{repo: clkRepository}
}

func (sch *storeClickHandler) HandleClick(ctx context.Context, click *entity.Click) <-chan *ClickProcessingResult {
	output := make(chan *ClickProcessingResult)

	go func(ctx context.Context, click *entity.Click) {
		defer close(output)
		output <- &ClickProcessingResult{
			Click: click,
			Err:   sch.repo.Save(ctx, click),
		}
	}(ctx, click)

	return output
}
