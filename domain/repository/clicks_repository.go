// Package repository contains a list of interfaces used by interactors.
package repository

import (
	"context"

	"github.com/lroman242/redirector/domain/entity"
)

// ClicksRepository interface describes clicks storage repository.
//
//go:generate mockgen -package=mocks -destination=mocks/mock_clicks_repository.go -source=domain/repository/clicks_repository.go ClicksRepository
type ClicksRepository interface {
	// Save function insert or update provided click to the storage.
	Save(ctx context.Context, click *entity.Click) error
}
