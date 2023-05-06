package dto

type RedirectResult struct {
	TargetURL string
	OutputCh  <-chan *ClickProcessingResult
}
