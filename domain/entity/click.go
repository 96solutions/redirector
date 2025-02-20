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

	TRKLink      *TrackingLink
	SourceID     string
	CampaignID   string
	AffiliateID  string
	AdvertiserID string
	IsParallel   bool

	LandingID string

	GCLID string

	UserAgent *valueobject.UserAgent
	Agent     string
	Platform  string //os
	Browser   string
	Device    string

	IP net.IP
	//Region      string
	CountryCode string
	//City        string

	P1 string
	P2 string
	P3 string
	P4 string

	CreatedAt time.Time
}
