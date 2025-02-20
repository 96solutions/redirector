// Package dto provides structures and functions for defining and handling data transfer objects.
// These DTOs are used to encapsulate data that is transferred between different layers of the application.
package dto

import (
	"net"
	"net/url"
)

// RedirectRequestData type describes set of data required for handling redirect request.
type RedirectRequestData struct {
	RequestID string
	Slug      string
	Params    map[string][]string
	Headers   map[string][]string
	UserAgent string
	IP        net.IP
	Protocol  string
	Referer   string
	URL       *url.URL
}

// GetParam is a helper function for convenient access to the request query params.
func (rrd *RedirectRequestData) GetParam(key string) []string {
	if val, ok := rrd.Params[key]; ok {
		return val
	}

	return make([]string, 0, 0)
}
