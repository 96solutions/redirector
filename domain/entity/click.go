package entity

import (
	"github.com/lroman242/redirector/domain/valueobject"
	"net"
)

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
