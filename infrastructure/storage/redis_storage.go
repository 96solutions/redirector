package storage

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/infrastructure/logger"
	"github.com/redis/go-redis/v9"
)

const (
	// trackingLinkKeyPrefix is used to create Redis keys for tracking links
	trackingLinkKeyPrefix = "trk:"
)

// RedisStorage implements repository.TrackingLinksRepositoryInterface using Redis
type RedisStorage struct {
	client *redis.Client
}

// NewRedisStorage creates a new RedisStorage instance
func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
	}
}

// FindTrackingLink retrieves a tracking link from Redis by slug
func (s *RedisStorage) FindTrackingLink(ctx context.Context, slug string) *entity.TrackingLink {
	// Get tracking link data
	data, err := s.client.Get(ctx, s.makeTrackingLinkKey(slug)).Bytes()
	if err != nil {
		if err != redis.Nil {
			slog.Error("failed to get tracking link from redis",
				"slug", slug,
				logger.ErrAttr(err),
			)
		}
		return nil
	}

	// Unmarshal tracking link with landing pages
	trkLink := new(entity.TrackingLink)
	if err := json.Unmarshal(data, trkLink); err != nil {
		slog.Error("failed to unmarshal tracking link",
			"slug", slug,
			logger.ErrAttr(err),
		)
		return nil
	}

	return trkLink
}

// makeTrackingLinkKey creates a Redis key for a tracking link
func (s *RedisStorage) makeTrackingLinkKey(slug string) string {
	return trackingLinkKeyPrefix + slug
}
