// Package service contains types which provide some specific functionality used for business logic.
package service

import "net"

//go:generate mockgen -package=mocks -destination=mocks/mock_ip_address_parser.go -source=ip_address_parser.go IPAddressParserInterface

// IPAddressParserInterface describes service that parses country code from the provided IP address.
type IPAddressParserInterface interface {
	// Parse function parses country code from the provided IP address.
	Parse(ip net.IP) (countryCode string, err error)
}
