package service

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoIP2 struct {
	db *geoip2.Reader
}

func NewGeoIP2(db *geoip2.Reader) *GeoIP2 {
	return &GeoIP2{db: db}
}

func (g *GeoIP2) Parse(ip net.IP) (string, error) {
	record, err := g.db.Country(ip)
	if err != nil {
		return "", err
	}

	return record.Country.IsoCode, nil
}
