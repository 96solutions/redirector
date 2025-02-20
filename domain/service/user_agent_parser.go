// Package service contains types which provide some specific functionality used for business logic.
package service

import (
	"github.com/lroman242/redirector/domain/valueobject"
)

//mockgen -package=mocks -destination=mocks/mock_user_agent_parser.go -source=domain/service/user_agent_parser.go UserAgentParser
//go:generate mockgen -package=mocks -destination=mocks/mock_user_agent_parser.go -source=user_agent_parser.go UserAgentParser

// UserAgentParserInterface interface describes service which parse data about used device from User-Agent header.
type UserAgentParserInterface interface {
	// Parse function parses data about used device from User-Agent header.
	Parse(userAgent string) (*valueobject.UserAgent, error)
}
