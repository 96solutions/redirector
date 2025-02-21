// Package dto provides structures and functions for defining and handling data transfer objects.
// These DTOs are used to encapsulate data that is transferred between different layers of the application.
package dto

import (
	"errors"
	"net"
	"net/url"
)

// RedirectRequestData contains all the information needed to process a redirect request,
// including request parameters, headers, user agent info, IP address, and protocol.
type RedirectRequestData struct {
	// RequestID uniquely identifies this redirect request
	RequestID string
	// Slug identifies the tracking link to use
	Slug string
	// Params contains URL query parameters
	Params map[string][]string
	// Headers contains HTTP request headers
	Headers map[string][]string
	// UserAgent is the raw User-Agent header string
	UserAgent string
	// IP is the client's IP address
	IP net.IP
	// Protocol is the request protocol (http/https)
	Protocol string
	// Referer contains the referring URL
	Referer string
	// URL contains the full request URL
	URL *url.URL
}

// GetParam is a helper function for convenient access to the request query params.
func (rrd *RedirectRequestData) GetParam(key string) []string {
	if val, ok := rrd.Params[key]; ok {
		return val
	}

	return make([]string, 0, 0)
}

// Validate function validates the redirect request data.
func (rrd *RedirectRequestData) Validate() error {
	if rrd.Slug == "" {
		return errors.New("slug is required")
	}
	if rrd.IP == nil {
		return errors.New("IP is required")
	}
	//TODO: Add more validation rules
	return nil
}
