package repository

import (
	"context"
	"github.com/lroman242/redirector/domain/entity"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_tracking_links_repository.go -source=domain/repository/tracking_links_repository.go TrackingLinksRepositoryInterface

type TrackingLinksRepositoryInterface interface {
	FindTrackingLink(ctx context.Context, slug string) *entity.TrackingLink
}
