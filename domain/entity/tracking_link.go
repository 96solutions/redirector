package entity

import "github.com/lroman242/redirector/domain/valueobject"

// TrackingLink type describes set of rules and requirements that should be used for handling redirect request.
type TrackingLink struct {
	Slug string

	IsActive bool

	AllowedProtocols []string
	AllowedGeos      []string
	AllowedDevices   []string

	IsCampaignOveraged           bool
	CampaignOverageRedirectRules *valueobject.RedirectRules

	IsCampaignActive              bool
	CampaignDisabledRedirectRules *valueobject.RedirectRules

	TargetURLTemplate string

	//AllowDeeplink bool
	//
	CampaignID   string
	AffiliateID  string
	AdvertiserID string
	SourceID     string
}
