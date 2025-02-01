// Package dto provides structures and functions for defining and handling data transfer objects.
// These DTOs are used to encapsulate data that is transferred between different layers of the application.
package dto

import "github.com/lroman242/redirector/domain/entity"

// ClickProcessingResult type describes output of interactor.ClickHandlerInterface.
type ClickProcessingResult struct {
	Click *entity.Click
	Err   error
}
