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
	IsCampaignOveraged bool
	// CampaignOveragedRedirectRulesID references rules for overaged campaigns
	CampaignOveragedRedirectRulesID int32
	// CampaignOverageRedirectRules contains redirect logic for overaged campaigns
	CampaignOverageRedirectRules *valueobject.RedirectRules

	// IsCampaignActive indicates if the campaign is currently active
	IsCampaignActive bool
	// CampaignActiveRedirectRulesID references rules for inactive campaigns
	CampaignActiveRedirectRulesID int32
	// CampaignDisabledRedirectRules contains redirect logic for inactive campaigns
	CampaignDisabledRedirectRules *valueobject.RedirectRules

	// AllowedGeos defines which geographic locations are permitted
	AllowedGeos AllowedListType
	// CampaignGeoRedirectRulesID references geo-specific redirect rules
	CampaignGeoRedirectRulesID int32
	// CampaignGeoRedirectRules contains redirect logic for geo restrictions
	CampaignGeoRedirectRules *valueobject.RedirectRules

	// AllowedDevices defines which device types are permitted
	AllowedDevices AllowedListType
	// CampaignDevicesRedirectRulesID references device-specific redirect rules
	CampaignDevicesRedirectRulesID int32
	// CampaignDevicesRedirectRules contains redirect logic for device restrictions
	CampaignDevicesRedirectRules *valueobject.RedirectRules

	// AllowedOS defines which operating systems are permitted
	AllowedOS AllowedListType
	// CampaignOSRedirectRulesID references OS-specific redirect rules
	CampaignOSRedirectRulesID int32
	// CampaignOSRedirectRules contains redirect logic for OS restrictions
	CampaignOSRedirectRules *valueobject.RedirectRules

	// TargetURLTemplate is the template for generating the final redirect URL
	TargetURLTemplate string

	// AllowDeeplink indicates if deeplink redirects are allowed
	AllowDeeplink bool

	// CampaignID identifies the campaign this link belongs to
	CampaignID string
	// AffiliateID identifies the affiliate
	AffiliateID string
	// AdvertiserID identifies the advertiser
	AdvertiserID string
	// SourceID identifies the traffic source
	SourceID string

	// LandingPages maps landing page IDs to their configurations
	LandingPages map[string]*LandingPage
}

// AllowedListType represents a map of allowed values where the key is the value name
// and the boolean indicates if it's allowed.
type AllowedListType map[string]bool

// Value returns the JSON-encoded representation.
func (a AllowedListType) Value() (driver.Value, error) {
	return json.Marshal(map[string]bool(a))
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

	*a = x
	return nil
}

// LandingPage contains information about a landing page associated with a tracking link.
type LandingPage struct {
	// ID uniquely identifies this landing page
	ID string
	// Title is the display name of the landing page
	Title string
	// PreviewURL is the URL where the landing page can be previewed
	PreviewURL string
	// TargetURL is the actual URL where traffic will be sent
	TargetURL string
}
