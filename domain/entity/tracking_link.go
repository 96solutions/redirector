// Package entity contains files which describe business objects.
package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/lroman242/redirector/domain/valueobject"
)

// TrackingLink type describes set of rules and requirements that should be used for handling redirect request.
type TrackingLink struct {
	Slug string

	IsActive bool

	AllowedProtocols                AllowedListType
	CampaignProtocolRedirectRulesID int32
	CampaignProtocolRedirectRules   *valueobject.RedirectRules

	IsCampaignOveraged              bool
	CampaignOveragedRedirectRulesID int32
	CampaignOverageRedirectRules    *valueobject.RedirectRules

	IsCampaignActive              bool
	CampaignActiveRedirectRulesID int32
	CampaignDisabledRedirectRules *valueobject.RedirectRules

	AllowedGeos                AllowedListType
	CampaignGeoRedirectRulesID int32
	CampaignGeoRedirectRules   *valueobject.RedirectRules

	AllowedDevices                 AllowedListType
	CampaignDevicesRedirectRulesID int32
	CampaignDevicesRedirectRules   *valueobject.RedirectRules

	AllowedOS                 AllowedListType
	CampaignOSRedirectRulesID int32
	CampaignOSRedirectRules   *valueobject.RedirectRules

	TargetURLTemplate string

	AllowDeeplink bool

	CampaignID   string
	AffiliateID  string
	AdvertiserID string
	SourceID     string

	LandingPages map[string]*LandingPage
}

// AllowedListType represents custom type used to convert list of redirection filters into JSON (JSONb).
type AllowedListType map[string]bool

// Value returns the JSON-encoded representation.
func (a AllowedListType) Value() (driver.Value, error) {
	x := make(map[string]bool)

	return json.Marshal(x)
}

// Scan decodes a JSON-encoded value.
func (a *AllowedListType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Unmarshal from json to map[string]bool
	x := make(map[string]bool)
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	return nil
}

// LandingPage type represents data related to tracking links.
type LandingPage struct {
	ID         string
	Title      string
	PreviewURL string
	TargetURL  string
}
