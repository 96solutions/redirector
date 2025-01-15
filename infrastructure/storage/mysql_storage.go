package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/lroman242/redirector/domain/entity"
)

const findTrackingLinkBySlugQuery = `select * from tracking_links where slug = ? limit 1`

type MySQLStorage struct {
	*sql.DB
}

func NewMySQLStorage(host, port, user, pass, database string) *MySQLStorage {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, database))
	if err != nil {
		panic(err)
	}

	return &MySQLStorage{
		db,
	}
}

func (s *MySQLStorage) FindTrackingLink(ctx context.Context, slug string) *entity.TrackingLink {
	stmt, err := s.DB.PrepareContext(ctx, findTrackingLinkBySlugQuery)
	if err != nil {
		log.Printf("an error occured while preparing statement, error: %s", err)
		return nil
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, slug)
	if err != nil {
		log.Printf("an error occured while executing statement, error: %s", err)
		return nil
	}
	if result.Err() != nil {
		log.Printf("an error occured while executing statement, error: %s", err)
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
		log.Printf("an error occured while scaning query result, error: %s", err)
		return nil
	}

	return trkLink
}
