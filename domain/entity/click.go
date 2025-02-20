// Package entity contains files which describe business objects.
package entity

import (
	"net"
	"time"

	"github.com/lroman242/redirector/domain/valueobject"
)

// Click type describes scope of data retrieved during redirect request via tracking link.
type Click struct {
	ID         string
	TargetURL  string
	Referer    string
	TrkURL     string
	Slug       string
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
