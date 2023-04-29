package service

import "net"

//go:generate mockgen -package=mocks -destination=mocks/mock_ip_address_parser.go -source=domain/service/ip_address_parser.go IpAddressParserInterface

// IpAddressParserInterface describes service that parses country code from the provided IP address.
type IpAddressParserInterface interface {
	// Parse function parses country code from the provided IP address.
	Parse(ip net.IP) (countryCode string, err error)
}
