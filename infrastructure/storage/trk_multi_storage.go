package storage

import (
	"context"

	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/domain/repository"
)

// MultiStorage implements repository.TrackingLinksRepositoryInterface by querying
// multiple storage backends sequentially until a result is found.
type MultiStorage struct {
	// storages is a slice of storage backends to query in order
	storages []repository.TrackingLinksRepositoryInterface
}

// NewMultiStorage creates a new MultiStorage instance with the provided storage backends.
// The order of storages determines the search order - first storage will be queried first.
func NewMultiStorage(storages []repository.TrackingLinksRepositoryInterface) *MultiStorage {
	return &MultiStorage{storages}
}

// FindTrackingLink searches for a tracking link in each storage sequentially.
// Returns the first found result or nil if not found in any storage.
func (ms *MultiStorage) FindTrackingLink(ctx context.Context, slug string) *entity.TrackingLink {
	// Check context before starting search
	if ctx.Err() != nil {
		return nil
	}

	// Search through storages in order
	for _, storage := range ms.storages {
		select {
		case <-ctx.Done():
			// Context cancelled, stop searching
			return nil
		default:
			// Try to find tracking link in current storage
			if result := storage.FindTrackingLink(ctx, slug); result != nil {
				return result
			}
			// Not found in current storage, continue to next one
		}
	}

	// Not found in any storage
	return nil
}
