package interactor

import (
	"context"
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
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
	TrackingLinkNotFoundError = errors.New("no tracking link was found by slug")
)

const (
	IPAddressToken    = "{ip}"
	ClickIDToken      = "{click_id}"
	UserAgentToken    = "{user_agent}"
	CampaignIDToken   = "{campaign_id}"
	AffiliateIDToken  = "{aff_id}"
	SourceIDToken     = "{source_id}"
	AdvertiserIDToken = "{advertiser_id}"
	DateToken         = "{date}"
	DateTimeToken     = "{date_time}"
	TimestampToken    = "{timestamp}"
	P1Token           = "{p1}"
	P2Token           = "{p2}"
	P3Token           = "{p3}"
	P4Token           = "{p4}"
	CountryCodeToken  = "{country_code}"
	RefererToken      = "{referer}"
	RandomStrToken    = "{random_str}"
	RandomIntToken    = "{random_int}"
	DeviceToken       = "{device}"
	PlatformToken     = "{platform}"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_redirect_interactor.go -source=domain/interactor/redirect_interactor.go RedirectInteractor
type RedirectInteractor interface {
	Redirect(ctx context.Context, slug string, requestData *dto.RedirectRequestData) (string, error)
}

type redirectInteractor struct {
	trackingLinksRepo repository.TrackingLinksRepositoryInterface
	ipAddressParser   service.IpAddressParserInterface
	userAgentParser   service.UserAgentParser
	tokenRegExp       *regexp.Regexp
}

func NewRedirectInteractor(trkRepo repository.TrackingLinksRepositoryInterface, ipAddressParser service.IpAddressParserInterface, userAgentParser service.UserAgentParser) RedirectInteractor {
	compiledRegExp, err := regexp.Compile(`{({)?(\w+)(})?}`)
	if err != nil {
		panic(err)
	}

	return &redirectInteractor{
		trackingLinksRepo: trkRepo,
		ipAddressParser:   ipAddressParser,
		userAgentParser:   userAgentParser,
		tokenRegExp:       compiledRegExp,
	}
}

func (r *redirectInteractor) Redirect(ctx context.Context, slug string, requestData *dto.RedirectRequestData) (string, error) {
	trackingLink := r.trackingLinksRepo.FindTrackingLink(slug)
	if trackingLink == nil {
		return "", TrackingLinkNotFoundError
	}

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

	//if deeplinkURL, ok := requestData.Params["deeplink"]; ok && trackingLink.AllowDeeplink {
	//	//TODO: handle deeplink
	//}

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
	targetURL := trackingLink.TargetURLTemplate

	tokens := r.tokenRegExp.FindAllString(targetURL, -1)
	for _, token := range tokens {
		switch token {
		case IPAddressToken:
			targetURL = strings.Replace(targetURL, token, requestData.IP.String(), -1)
		case ClickIDToken:
			targetURL = strings.Replace(targetURL, token, requestData.RequestID, -1)
		case UserAgentToken:
			targetURL = strings.Replace(targetURL, token, requestData.UserAgent, -1)
		case CampaignIDToken:
			targetURL = strings.Replace(targetURL, token, trackingLink.CampaignID, -1)
		case AffiliateIDToken:
			targetURL = strings.Replace(targetURL, token, trackingLink.AffiliateID, -1)
		case SourceIDToken:
			targetURL = strings.Replace(targetURL, token, trackingLink.SourceID, -1)
		case AdvertiserIDToken:
			targetURL = strings.Replace(targetURL, token, trackingLink.AdvertiserID, -1)
		case DateToken:
			targetURL = strings.Replace(targetURL, token, time.Now().Format("2006-02-01"), 1)
		case DateTimeToken:
			targetURL = strings.Replace(targetURL, token, time.Now().Format("2006-01-02T15:04:05"), 1)
		case TimestampToken:
			targetURL = strings.Replace(targetURL, token, strconv.FormatInt(time.Now().Unix(), 10), 1)
		case P1Token:
			values := requestData.GetParam("p1")
			targetURL = strings.Replace(targetURL, token, strings.Join(values, ","), -1)
		case P2Token:
			values := requestData.GetParam("p2")
			targetURL = strings.Replace(targetURL, token, strings.Join(values, ","), -1)
		case P3Token:
			values := requestData.GetParam("p3")
			targetURL = strings.Replace(targetURL, token, strings.Join(values, ","), -1)
		case P4Token:
			values := requestData.GetParam("p4")
			targetURL = strings.Replace(targetURL, token, strings.Join(values, ","), -1)
		case CountryCodeToken:
			targetURL = strings.Replace(targetURL, token, countryCode, -1)
		case RefererToken:
			targetURL = strings.Replace(targetURL, token, requestData.Referer, -1)
		case RandomStrToken:
			targetURL = strings.Replace(targetURL, token, randString(32), 1)
		case RandomIntToken:
			targetURL = strings.Replace(targetURL, token, strconv.Itoa(rand.Intn(99999999-10000)+10000), 1)
		case DeviceToken:
			targetURL = strings.Replace(targetURL, token, ua.Device, -1)
		case PlatformToken:
			targetURL = strings.Replace(targetURL, token, ua.Platform, -1)
		//TODO: add new tokens

		//replace undefined tokens with empty string
		default:
			targetURL = strings.Replace(targetURL, token, "", -1)
		}
	}

	return targetURL
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randString function generates random string of n-length
func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
