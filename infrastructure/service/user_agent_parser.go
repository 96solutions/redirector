package service

import (
	"errors"

	"github.com/lroman242/redirector/domain/service"
	"github.com/lroman242/redirector/domain/valueobject"
	"github.com/ua-parser/uap-go/uaparser"
)

// ErrEmptyUserAgent is returned when the provided User-Agent string is empty.
var ErrEmptyUserAgent = errors.New("provided empty user agent")

// UserAgentParser implements service.UserAgentParserInterface.
type UserAgentParser struct {
	*uaparser.Parser
}

// NewUserAgentParser func creates new instance of UserAgentParser.
func NewUserAgentParser() service.UserAgentParserInterface {
	return &UserAgentParser{
		uaparser.NewFromSaved(),
	}
}

// Parse function parses data about used device from User-Agent header.
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
