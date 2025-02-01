// Package dto provides structures and functions for defining and handling data transfer objects.
// These DTOs are used to encapsulate data that is transferred between different layers of the application.
package dto

// RedirectResult type describes output of interactor.RedirectInteractor Redirect function.
type RedirectResult struct {
	TargetURL string
	OutputCh  <-chan *ClickProcessingResult
}
