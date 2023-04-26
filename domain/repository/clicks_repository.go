package repository

import (
	"context"

	"github.com/lroman242/redirector/domain/entity"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_clicks_repository.go -source=domain/repository/clicks_repository.go ClicksRepository
type ClicksRepository interface {
	Save(ctx context.Context, click *entity.Click) error
}
