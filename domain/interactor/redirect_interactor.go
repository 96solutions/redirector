package interactor

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/lroman242/redirector/domain/dto"
	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/domain/repository"
	"github.com/lroman242/redirector/domain/service"
	"github.com/lroman242/redirector/domain/valueobject"
)

var (
	UnsupportedProtocolError  = errors.New("protocol is not allowed for that tracking link")
	UnsupportedGeoError       = errors.New("visitor geo is not allowed for that tracking link")
	UnsupportedDeviceError    = errors.New("visitor device is not allowed for that tracking link")
	InvalidRedirectTypeError  = errors.New("invalid redirect type is stored in tracking link redirect rules")
	BlockRedirectError        = errors.New("redirect should be blocked")
	TrackingLinkDisabledError = errors.New("used tracking link is disabled")
)

//go:generate mockgen -package=mocks -destination=mocks/mock_redirect_interactor.go -source=domain/interactor/redirect_interactor.go RedirectInteractor
type RedirectInteractor interface {
	Redirect(ctx context.Context, slug string, requestData *dto.RedirectRequestData) (string, error)
}

type redirectInteractor struct {
	trackingLinksRepo repository.TrackingLinksRepositoryInterface
	ipAddressParser   service.IpAddressParserInterface
	userAgentParser   service.UserAgentParser
}

func NewRedirectInteractor(trkRepo repository.TrackingLinksRepositoryInterface, ipAddressParser service.IpAddressParserInterface, userAgentParser service.UserAgentParser) RedirectInteractor {
	return &redirectInteractor{
		trackingLinksRepo: trkRepo,
		ipAddressParser:   ipAddressParser,
		userAgentParser:   userAgentParser,
	}
}

func (r *redirectInteractor) Redirect(ctx context.Context, slug string, requestData *dto.RedirectRequestData) (string, error) {
	trackingLink := r.trackingLinksRepo.FindTrackingLink(slug)
	//TODO: handle tracking link not found case

	if !trackingLink.IsActive {
		return "", TrackingLinkDisabledError
	}

	if len(trackingLink.AllowedProtocols) > 0 && !contains(requestData.Protocol, trackingLink.AllowedProtocols) {
		return "", UnsupportedProtocolError
	}

	countryCode, err := r.ipAddressParser.Parse(requestData.IP)
	if err != nil {
		//TODO: log error
	}
	if len(trackingLink.AllowedGeos) > 0 && !contains(countryCode, trackingLink.AllowedGeos) {
		return "", UnsupportedGeoError
	}

	ua, err := r.userAgentParser.Parse(requestData.UserAgent)
	if err != nil {
		//TODO: log error
	}
	if len(trackingLink.AllowedDevices) > 0 && !contains(ua.Device, trackingLink.AllowedDevices) {
		return "", UnsupportedDeviceError
	}

	if trackingLink.IsCampaignOveraged {
		return r.handleRedirectRules(trackingLink.CampaignOverageRedirectRules, ctx, requestData)
	}

	if !trackingLink.IsCampaignActive {
		return r.handleRedirectRules(trackingLink.CampaignDisabledRedirectRules, ctx, requestData)
	}

	targetURL := r.renderTokens(trackingLink, requestData, ua, countryCode)

	//TODO: implement pipe: service -> click registration -> [kafka producer, clickhouse insert]

	return targetURL, nil
}

func (r *redirectInteractor) handleRedirectRules(rr *valueobject.RedirectRules, ctx context.Context, requestData *dto.RedirectRequestData) (string, error) {
	switch rr.RedirectType {
	case valueobject.LinkRedirectType:
		return rr.RedirectURL, nil
	case valueobject.SlugRedirectType:
		return r.Redirect(ctx, rr.RedirectSlug, requestData)
	case valueobject.SmartSlugRedirectType:
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		newSlug := rr.RedirectSmartSlug[rnd.Intn(len(rr.RedirectSmartSlug))]
		return r.Redirect(ctx, newSlug, requestData)
	case valueobject.NoRedirectType:
		return "", BlockRedirectError
	default:
		return "", InvalidRedirectTypeError
	}
}

func (r *redirectInteractor) renderTokens(trackingLink *entity.TrackingLink, requestData *dto.RedirectRequestData, ua *valueobject.UserAgent, countryCode string) string {
	//TODO: render tokens for target URL

	return trackingLink.TargetURLTemplate
}

// contains checks if a string is present in a slice.
func contains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}
