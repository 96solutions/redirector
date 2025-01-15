package service

import (
	"errors"

	"github.com/lroman242/redirector/domain/service"
	"github.com/lroman242/redirector/domain/valueobject"
	"github.com/ua-parser/uap-go/uaparser"
)

var EmptyUserAgentError = errors.New("provided empty user agent")

type UserAgentParser struct {
	*uaparser.Parser
}

func NewUserAgentParser() service.UserAgentParser {
	return &UserAgentParser{
		uaparser.NewFromSaved(),
	}
}

func (p *UserAgentParser) Parse(userAgent string) (*valueobject.UserAgent, error) {
	ua := new(valueobject.UserAgent)

	if len(userAgent) == 0 {
		return ua, EmptyUserAgentError
	}

	client := p.Parser.Parse(userAgent)

	ua.SrcString = userAgent
	ua.Device = client.Device.Family
	ua.Platform = client.Os.Family
	ua.Browser = client.UserAgent.Family

	return ua, nil
}
