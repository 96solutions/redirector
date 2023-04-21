package repository

import "github.com/lroman242/redirector/domain/entity"

//go:generate mockgen -package=mocks -destination=mocks/mock_tracking_link_repository.go -source=domain/repository/tracking_links.go TrackingLinksRepositoryInterface

type TrackingLinksRepositoryInterface interface {
	FindTrackingLink(slug string) *entity.TrackingLink
}
