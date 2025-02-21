// Package entity contains files which describe business objects.
package entity

import (
	"net"
	"time"

	"github.com/lroman242/redirector/domain/valueobject"
)

// Click represents a tracking event generated when a redirect occurs.
// It stores information about the redirect request including user data,
// campaign information, and tracking parameters.
type Click struct {
	// ID uniquely identifies this click event
	ID string
	// TargetURL is the final URL where user was redirected
	TargetURL string
	// Referer contains the referring URL
	Referer string
	// TrkURL is the tracking URL that was accessed
	TrkURL string
	// Slug identifies the tracking link used
	Slug string
	// ParentSlug identifies the parent tracking link if this was a chained redirect
	ParentSlug string

	// TRKLink is the tracking link that was used
	TRKLink *TrackingLink
	// SourceID is the source ID
	SourceID string
	// CampaignID is the campaign ID
	CampaignID string
	// AffiliateID is the affiliate ID
	AffiliateID string
	// AdvertiserID is the advertiser ID
	AdvertiserID string
	// IsParallel indicates if this is a parallel redirect
	IsParallel bool

	// LandingID is the landing page ID
	LandingID string
	// GCLID is the Google Click Identifier
	GCLID string

	// UserAgent is the user agent
	UserAgent *valueobject.UserAgent
	Agent     string
	Platform  string //os
	Browser   string
	Device    string

	// IP is the IP address
	IP net.IP
	// CountryCode is the country code
	CountryCode string

	// P1-P4 are custom parameters
	P1 string
	P2 string
	P3 string
	P4 string

	// CreatedAt is the creation time
	CreatedAt time.Time
}
