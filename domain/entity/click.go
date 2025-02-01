// Package entity contains files which describe business objects.
package entity

import (
	"net"

	"github.com/lroman242/redirector/domain/valueobject"
)

// Click type describes scope of data retrieved during redirect request via tracking link.
type Click struct {
	ID          string
	TargetURL   string
	TRKLink     *TrackingLink
	UserAgent   *valueobject.UserAgent
	CountryCode string
	IP          net.IP
	Referer     string
	P1          string
	P2          string
	P3          string
	P4          string
	//TODO: add more field
}
