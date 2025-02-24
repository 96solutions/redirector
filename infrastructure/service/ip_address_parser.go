package service

import (
	"errors"
	"fmt"
	"net"

	"github.com/lroman242/redirector/domain/service"
	"github.com/oschwald/geoip2-golang"
)

// GeoIP2 implements service.IPAddressParserInterface, allows to find Country name from ip address.
type GeoIP2 struct {
	db *geoip2.Reader
}

// NewGeoIP2 func creates new instance of GeoIP2.
func NewGeoIP2(db *geoip2.Reader) service.IPAddressParserInterface {
	return &GeoIP2{db: db}
}

// Parse function parses country code from the provided IP address.
func (g *GeoIP2) Parse(ip net.IP) (string, error) {
	if ip == nil {
		return "", errors.New("ip address cannot be nil")
	}

	record, err := g.db.Country(ip)
	if err != nil {
		return "", fmt.Errorf("failed to get country from IP: %w", err)
	}

	return record.Country.IsoCode, nil
}

// Close function closes connection to the geo ip database.
func (g *GeoIP2) Close() error {
	return g.db.Close()
}
