package storage

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/domain/valueobject"
	"github.com/lroman242/redirector/infrastructure/logger"
)

// findTrackingLinkBySlugQuery selects all tracking link data including redirect rules.
const findTrackingLinkBySlugQuery = `SELECT
    t.slug,
    t.active,
    t.allowed_protocols,
    t.allowed_geos,
    t.allowed_devices,
    t.allowed_os,
    t.campaign_overaged,
    t.campaign_overaged_redirect_rules_id,
    ovr.redirect_type as campaign_overaged_redirect_type,
    ovr.redirect_slug as campaign_overaged_redirect_slug,
    ovr.redirect_url as campaign_overaged_redirect_url,
    ovr.redirect_smart_slug as campaign_overaged_redirect_smart_slug,
    t.campaign_active,
    t.campaign_active_redirect_rules_id,
    dr.redirect_type as campaign_disabled_redirect_type,
    dr.redirect_slug as campaign_disabled_redirect_slug,
    dr.redirect_url as campaign_disabled_redirect_url,
    dr.redirect_smart_slug as campaign_disabled_redirect_smart_slug,
    t.campaign_protocol_redirect_rules_id,
    pr.redirect_type as protocol_redirect_type,
    pr.redirect_slug as protocol_redirect_slug,
    pr.redirect_url as protocol_redirect_url,
    pr.redirect_smart_slug as protocol_redirect_smart_slug,
    t.campaign_geo_redirect_rules_id,
    gr.redirect_type as geo_redirect_type,
    gr.redirect_slug as geo_redirect_slug,
    gr.redirect_url as geo_redirect_url,
    gr.redirect_smart_slug as geo_redirect_smart_slug,
    t.campaign_devices_redirect_rules_id,
    devr.redirect_type as devices_redirect_type,
    devr.redirect_slug as devices_redirect_slug,
    devr.redirect_url as devices_redirect_url,
    devr.redirect_smart_slug as devices_redirect_smart_slug,
    t.campaign_os_redirect_rules_id,
    osr.redirect_type as os_redirect_type,
    osr.redirect_slug as os_redirect_slug,
    osr.redirect_url as os_redirect_url,
    osr.redirect_smart_slug as os_redirect_smart_slug,
    t.target_url_template,
    t.allow_deeplink,
    t.campaign_id,
    t.affiliate_id,
    t.advertiser_id,
    t.source_id
FROM tracking_links t
LEFT JOIN redirect_rules ovr ON ovr.id = t.campaign_overaged_redirect_rules_id
LEFT JOIN redirect_rules dr ON dr.id = t.campaign_active_redirect_rules_id
LEFT JOIN redirect_rules pr ON pr.id = t.campaign_protocol_redirect_rules_id
LEFT JOIN redirect_rules gr ON gr.id = t.campaign_geo_redirect_rules_id
LEFT JOIN redirect_rules devr ON devr.id = t.campaign_devices_redirect_rules_id
LEFT JOIN redirect_rules osr ON osr.id = t.campaign_os_redirect_rules_id
WHERE t.slug = $1
LIMIT 1`

// SQLStorage implements repository.TrackingLinksRepositoryInterface.
type SQLStorage struct {
	*sql.DB
}

// NewSQLStorage creates a new SQLStorage instance.
func NewSQLStorage(dbConnection *sql.DB) *SQLStorage {
	return &SQLStorage{
		dbConnection,
	}
}

// FindTrackingLink retrieves a tracking link and its associated redirect rules by slug.
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
	defer result.Close()

	if !result.Next() {
		return nil
	}

	trkLink := new(entity.TrackingLink)
	trkLink.CampaignOverageRedirectRules = new(valueobject.RedirectRules)
	trkLink.CampaignDisabledRedirectRules = new(valueobject.RedirectRules)
	trkLink.CampaignProtocolRedirectRules = new(valueobject.RedirectRules)
	trkLink.CampaignGeoRedirectRules = new(valueobject.RedirectRules)
	trkLink.CampaignDevicesRedirectRules = new(valueobject.RedirectRules)
	trkLink.CampaignOSRedirectRules = new(valueobject.RedirectRules)

	err = result.Scan(
		&trkLink.Slug,
		&trkLink.IsActive,
		&trkLink.AllowedProtocols,
		&trkLink.AllowedGeos,
		&trkLink.AllowedDevices,
		&trkLink.AllowedOS,

		&trkLink.IsCampaignOveraged,
		&trkLink.CampaignOveragedRedirectRulesID,
		&trkLink.CampaignOverageRedirectRules.RedirectType,
		&trkLink.CampaignOverageRedirectRules.RedirectSlug,
		&trkLink.CampaignOverageRedirectRules.RedirectURL,
		&trkLink.CampaignOverageRedirectRules.RedirectSmartSlug,

		&trkLink.IsCampaignActive,
		&trkLink.CampaignActiveRedirectRulesID,
		&trkLink.CampaignDisabledRedirectRules.RedirectType,
		&trkLink.CampaignDisabledRedirectRules.RedirectSlug,
		&trkLink.CampaignDisabledRedirectRules.RedirectURL,
		&trkLink.CampaignDisabledRedirectRules.RedirectSmartSlug,

		&trkLink.CampaignProtocolRedirectRulesID,
		&trkLink.CampaignProtocolRedirectRules.RedirectType,
		&trkLink.CampaignProtocolRedirectRules.RedirectSlug,
		&trkLink.CampaignProtocolRedirectRules.RedirectURL,
		&trkLink.CampaignProtocolRedirectRules.RedirectSmartSlug,

		&trkLink.CampaignGeoRedirectRulesID,
		&trkLink.CampaignGeoRedirectRules.RedirectType,
		&trkLink.CampaignGeoRedirectRules.RedirectSlug,
		&trkLink.CampaignGeoRedirectRules.RedirectURL,
		&trkLink.CampaignGeoRedirectRules.RedirectSmartSlug,

		&trkLink.CampaignDevicesRedirectRulesID,
		&trkLink.CampaignDevicesRedirectRules.RedirectType,
		&trkLink.CampaignDevicesRedirectRules.RedirectSlug,
		&trkLink.CampaignDevicesRedirectRules.RedirectURL,
		&trkLink.CampaignDevicesRedirectRules.RedirectSmartSlug,

		&trkLink.CampaignOSRedirectRulesID,
		&trkLink.CampaignOSRedirectRules.RedirectType,
		&trkLink.CampaignOSRedirectRules.RedirectSlug,
		&trkLink.CampaignOSRedirectRules.RedirectURL,
		&trkLink.CampaignOSRedirectRules.RedirectSmartSlug,

		&trkLink.TargetURLTemplate,
		&trkLink.AllowDeeplink,
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
