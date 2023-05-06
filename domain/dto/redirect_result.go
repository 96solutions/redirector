package dto

// RedirectResult type describes output of interactor.RedirectInteractor Redirect function.
type RedirectResult struct {
	TargetURL string
	OutputCh  <-chan *ClickProcessingResult
}
