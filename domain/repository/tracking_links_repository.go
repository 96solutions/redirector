package repository

import (
	"context"

	"github.com/lroman242/redirector/domain/entity"
)

// TrackingLinksRepositoryInterface interface describes clicks storage repository.
//
//go:generate mockgen -package=mocks -destination=mocks/mock_tracking_links_repository.go -source=domain/repository/tracking_links_repository.go TrackingLinksRepositoryInterface
type TrackingLinksRepositoryInterface interface {
	// FindTrackingLink function retrieves entity.TrackingLink record from the storage by slug.
	FindTrackingLink(ctx context.Context, slug string) *entity.TrackingLink
}
