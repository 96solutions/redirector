package storage

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/domain/valueobject"
	"github.com/lroman242/redirector/infrastructure/logger"
)

const findTrackingLinkBySlugQuery = `SELECT  slug, active, allowed_protocols, allowed_geos, allowed_devices, campaign_overaged, campaign_overaged_redirect_rules_id, campaign_active, campaign_active_redirect_rules_id, target_url_template, campaign_id, publisher_id, advertiser_id, source_id FROM tracking_links WHERE slug = $1 LIMIT 1`

type SQLStorage struct {
	*sql.DB
}

func NewSQLStorage(dbConnection *sql.DB) *SQLStorage {
	return &SQLStorage{
		dbConnection,
	}
}

func (s *SQLStorage) FindTrackingLink(ctx context.Context, slug string) *entity.TrackingLink {
	stmt, err := s.DB.PrepareContext(ctx, findTrackingLinkBySlugQuery)
	if err != nil {
		slog.Error("an error occurred while preparing statement", logger.ErrAttr(err))
		return nil
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, slug)
	if err != nil {
		slog.Error("an error occurred while executing statement", logger.ErrAttr(err))
		return nil
	}
	if result.Err() != nil {
		slog.Error("an error occurred while executing statement", logger.ErrAttr(err))
		return nil
	}
	result.Next()

	trkLink := new(entity.TrackingLink)
	trkLink.CampaignOverageRedirectRules = new(valueobject.RedirectRules)
	trkLink.CampaignDisabledRedirectRules = new(valueobject.RedirectRules)

	err = result.Scan(
		&trkLink.Slug,
		&trkLink.IsActive,
		&trkLink.AllowedProtocols,
		&trkLink.AllowedGeos,
		&trkLink.AllowedDevices,

		&trkLink.IsCampaignOveraged,
		&trkLink.CampaignOveragedRedirectRulesID,
		//&trkLink.CampaignOverageRedirectRules.RedirectType,
		//&trkLink.CampaignOverageRedirectRules.RedirectSlug,
		//&trkLink.CampaignOverageRedirectRules.RedirectURL,
		//&trkLink.CampaignOverageRedirectRules.RedirectSmartSlug,

		&trkLink.IsCampaignActive,
		&trkLink.CampaignActiveRedirectRulesID,
		//&trkLink.CampaignDisabledRedirectRules.RedirectType,
		//&trkLink.CampaignDisabledRedirectRules.RedirectSlug,
		//&trkLink.CampaignDisabledRedirectRules.RedirectURL,
		//&trkLink.CampaignDisabledRedirectRules.RedirectSmartSlug,

		&trkLink.TargetURLTemplate,

		&trkLink.CampaignID,
		&trkLink.AffiliateID,
		&trkLink.AdvertiserID,
		&trkLink.SourceID,
	)

	if err != nil {
		slog.Error("an error occurred while scanning query result", logger.ErrAttr(err))
		return nil
	}

	return trkLink
}
