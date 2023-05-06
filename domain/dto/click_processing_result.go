package dto

import "github.com/lroman242/redirector/domain/entity"

// ClickProcessingResult type describes output of interactor.ClickHandlerInterface
type ClickProcessingResult struct {
	Click *entity.Click
	Err   error
}
