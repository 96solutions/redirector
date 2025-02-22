// Package service contains types which provide some specific functionality used for business logic.
package service

import (
	"github.com/lroman242/redirector/domain/valueobject"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_user_agent_parser.go -source=user_agent_parser.go UserAgentParser

// UserAgentParserInterface provides functionality to parse User-Agent strings
// and extract device, platform and browser information.
type UserAgentParserInterface interface {
	// Parse analyzes a User-Agent string and returns structured information about
	// the client's device, platform, and browser. Returns an error if parsing fails.
	Parse(userAgent string) (*valueobject.UserAgent, error)
}
