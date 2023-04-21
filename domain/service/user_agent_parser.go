package service

import (
	"github.com/lroman242/redirector/domain/valueobject"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_user_agent_parser.go -source=domain/service/user_agent_parser.go UserAgentParser

type UserAgentParser interface {
	Parse(userAgent string) (*valueobject.UserAgent, error)
}
