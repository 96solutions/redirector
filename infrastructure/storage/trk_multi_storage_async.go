package storage

import (
	"context"
	"sync"

	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/domain/repository"
)

// MultiStorageAsync implements repository.TrackingLinksRepositoryInterface by querying
// multiple storage backends concurrently and returning the first successful result.
type MultiStorageAsync struct {
	// storages is a slice of storage backends to query
	storages []repository.TrackingLinksRepositoryInterface
}

// NewMultiStorageAsync creates a new MultiStorageAsync instance with the provided storage backends.
// The order of storages determines the query priority.
func NewMultiStorageAsync(storages []repository.TrackingLinksRepositoryInterface) *MultiStorageAsync {
	return &MultiStorageAsync{storages}
}

// FindTrackingLink queries all storage backends concurrently and returns the first
// successful result. If no storage returns a result, returns nil.
func (ms *MultiStorageAsync) FindTrackingLink(ctx context.Context, slug string) *entity.TrackingLink {
	// Create channel for results with buffer size equal to number of storages
	results := make(chan *entity.TrackingLink, len(ms.storages))
	// Create WaitGroup to track goroutines
	var wg sync.WaitGroup

	// Query each storage concurrently
	for _, storage := range ms.storages {
		wg.Add(1)
		go func(s repository.TrackingLinksRepositoryInterface) {
			defer wg.Done()
			// Query the storage
			if result := s.FindTrackingLink(ctx, slug); result != nil {
				select {
				case results <- result:
					// Result sent successfully
				case <-ctx.Done():
					// Context cancelled
					return
				}
			}
		}(storage)
	}

	// Close results channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Return first non-nil result or nil if none found
	select {
	case result := <-results:
		return result
	case <-ctx.Done():
		return nil
	}
}
