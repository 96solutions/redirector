package service

import "net"

//go:generate mockgen -package=mocks -destination=mocks/mock_ip_address_parser.go -source=domain/service/ip_address_parser.go IpAddressParserInterface

type IpAddressParserInterface interface {
	Parse(ip net.IP) (countryCode string, err error)
}
