package storage

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/lroman242/redirector/domain/entity"
)

const findTrackingLinkBySlugQuery = `select * from tracking_links where slug = ? limit 1`

type MySQLStorage struct {
	*sql.DB
}

func NewMySQLStorage(dbConnection *sql.DB) *MySQLStorage {
	return &MySQLStorage{
		dbConnection,
	}
}

func (s *MySQLStorage) FindTrackingLink(ctx context.Context, slug string) *entity.TrackingLink {
	stmt, err := s.DB.PrepareContext(ctx, findTrackingLinkBySlugQuery)
	if err != nil {
		slog.Error("an error occurred while preparing statement", "error", err)
		return nil
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, slug)
	if err != nil {
		slog.Error("an error occurred while executing statement", "error", err)
		return nil
	}
	if result.Err() != nil {
		slog.Error("an error occurred while executing statement", "error", err)
		return nil
	}

	trkLink := new(entity.TrackingLink)
	err = result.Scan(
		&trkLink.Slug,
		&trkLink.IsActive,
		&trkLink.AllowedProtocols,
		&trkLink.AllowedGeos,
		&trkLink.AllowedDevices,

		&trkLink.IsCampaignOveraged,
		&trkLink.CampaignOverageRedirectRules.RedirectType,
		&trkLink.CampaignOverageRedirectRules.RedirectSlug,
		&trkLink.CampaignOverageRedirectRules.RedirectURL,
		&trkLink.CampaignOverageRedirectRules.RedirectSmartSlug,

		&trkLink.IsCampaignActive,
		&trkLink.CampaignDisabledRedirectRules.RedirectType,
		&trkLink.CampaignDisabledRedirectRules.RedirectSlug,
		&trkLink.CampaignDisabledRedirectRules.RedirectURL,
		&trkLink.CampaignDisabledRedirectRules.RedirectSmartSlug,

		&trkLink.TargetURLTemplate,

		&trkLink.CampaignID,
		&trkLink.AffiliateID,
		&trkLink.AdvertiserID,
		&trkLink.SourceID,
	)

	if err != nil {
		slog.Error("an error occurred while scanning query result", "error", err)
		return nil
	}

	return trkLink
}
