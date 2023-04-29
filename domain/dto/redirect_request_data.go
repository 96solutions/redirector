package dto

import (
	"net"
)

// RedirectRequestData type describes set of data required for handling redirect request.
type RedirectRequestData struct {
	Params    map[string][]string
	Headers   map[string]string
	UserAgent string
	IP        net.IP
	Protocol  string
	Referer   string

	RequestID string
}

// GetParam is a helper function for convenient access to the request query params.
func (rrd *RedirectRequestData) GetParam(key string) []string {
	if val, ok := rrd.Params[key]; ok {
		return val
	}

	return make([]string, 0, 0)
}
