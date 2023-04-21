package dto

import "net"

// RedirectRequestData type describes set of data required for handling redirect request.
type RedirectRequestData struct {
	Params    map[string][]string
	Headers   map[string]string
	UserAgent string
	IP        net.IP
	Protocol  string
}
