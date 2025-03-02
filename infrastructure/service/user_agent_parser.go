// Package service provides implementations of domain service interfaces and additional
// service wrappers for metrics, caching, and other cross-cutting concerns.
package service

import (
	"errors"

	"github.com/lroman242/redirector/domain/service"
	"github.com/lroman242/redirector/domain/valueobject"
	"github.com/ua-parser/uap-go/uaparser"
)

// ErrEmptyUserAgent is returned when the provided User-Agent string is empty.
var ErrEmptyUserAgent = errors.New("provided empty user agent")

// UserAgentParser implements service.UserAgentParserInterface using the ua-parser library.
// It provides functionality to parse User-Agent strings and extract device, platform,
// and browser information.
type UserAgentParser struct {
	// Parser is the underlying ua-parser implementation
	*uaparser.Parser
}

// NewUserAgentParser creates a new UserAgentParser instance using the default
// ua-parser patterns database.
func NewUserAgentParser() service.UserAgentParserInterface {
	return &UserAgentParser{
		uaparser.NewFromSaved(),
	}
}

// Parse analyzes a User-Agent string and returns structured information about
// the client's device, platform, and browser. Returns ErrEmptyUserAgent if
// the provided string is empty.
func (p *UserAgentParser) Parse(userAgent string) (*valueobject.UserAgent, error) {
	ua := new(valueobject.UserAgent)

	if len(userAgent) == 0 {
		return ua, ErrEmptyUserAgent
	}

	client := p.Parser.Parse(userAgent)

	ua.SrcString = userAgent
	ua.Device = client.Device.Family
	ua.Platform = client.Os.Family
	ua.Browser = client.UserAgent.Family

	return ua, nil
}
