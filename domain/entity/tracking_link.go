// Package entity contains files which describe business objects.
package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/lroman242/redirector/domain/valueobject"
)

// TrackingLink defines the rules and configuration for handling redirect requests.
// It includes validation rules for protocols, geos, and devices, as well as
// redirect rules for different campaign states (overaged, disabled, etc).
type TrackingLink struct {
	// Slug uniquely identifies this tracking link
	Slug string

	// IsActive indicates if this tracking link is enabled
	IsActive bool

	// AllowedProtocols defines which protocols (http/https) are permitted
	AllowedProtocols AllowedListType

	// CampaignProtocolRedirectRulesID references protocol-specific redirect rules
	CampaignProtocolRedirectRulesID int32
	// CampaignProtocolRedirectRules contains protocol-specific redirect logic
	CampaignProtocolRedirectRules *valueobject.RedirectRules

	// IsCampaignOveraged indicates if campaign limits have been exceeded
	IsCampaignOveraged              bool
	CampaignOveragedRedirectRulesID int32
	CampaignOverageRedirectRules    *valueobject.RedirectRules

	// IsCampaignActive indicates if the campaign is active
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

// AllowedListType represents a map of allowed values where the key is the value name
// and the boolean indicates if it's allowed
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

// LandingPage contains information about a landing page associated with a tracking link
type LandingPage struct {
	ID         string
	Title      string
	PreviewURL string
	TargetURL  string
}
