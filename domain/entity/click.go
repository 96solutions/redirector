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
	// SourceID identifies the traffic source
	SourceID string
	// CampaignID identifies the campaign
	CampaignID string
	// AffiliateID identifies the affiliate
	AffiliateID string
	// AdvertiserID identifies the advertiser
	AdvertiserID string
	// IsParallel indicates if this is a parallel redirect
	IsParallel bool

	// LandingID identifies the landing page
	LandingID string
	// GCLID is the Google Click Identifier
	GCLID string

	// UserAgent contains parsed user agent information
	UserAgent *valueobject.UserAgent
	// Agent is the raw user agent string
	Agent string
	// Platform is the operating system
	Platform string
	// Browser is the web browser name
	Browser string
	// Device is the device type
	Device string

	// IP is the visitor's IP address
	IP net.IP
	// CountryCode is the visitor's country code
	CountryCode string

	// P1-P4 are custom tracking parameters
	P1 string
	P2 string
	P3 string
	P4 string

	// CreatedAt is when this click was recorded
	CreatedAt time.Time
}
