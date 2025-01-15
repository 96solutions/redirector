// Package interactor contains all use-case interactors preformed by the application.
package interactor

import (
	"context"
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
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

// RedirectInteractor interface describes the service to handle requests and return the following target URL to redirect to.
type RedirectInteractor interface {
	Redirect(ctx context.Context, slug string, requestData *dto.RedirectRequestData) (*dto.RedirectResult, error)
}

type redirectInteractor struct {
	log                     service.Logger
	trackingLinksRepository repository.TrackingLinksRepositoryInterface
	ipAddressParser         service.IpAddressParserInterface
	userAgentParser         service.UserAgentParser
	tokenRegExp             *regexp.Regexp
	clickHandlers           []ClickHandlerInterface
}

// NewRedirectInteractor function creates RedirectInteractor implementation.
func NewRedirectInteractor(
	logger service.Logger,
	trkRepo repository.TrackingLinksRepositoryInterface,
	ipAddressParser service.IpAddressParserInterface,
	userAgentParser service.UserAgentParser,
	clickHandlers []ClickHandlerInterface,
) RedirectInteractor {
	compiledRegExp := regexp.MustCompile(`{({)?(\w+)(})?}`)

	return &redirectInteractor{
		log:                     logger,
		trackingLinksRepository: trkRepo,
		ipAddressParser:         ipAddressParser,
		userAgentParser:         userAgentParser,
		tokenRegExp:             compiledRegExp,
		clickHandlers:           clickHandlers,
	}
}

// Redirect function handles requests and returns the target URL to redirect traffic to.
func (r *redirectInteractor) Redirect(
	ctx context.Context,
	slug string,
	requestData *dto.RedirectRequestData,
) (*dto.RedirectResult, error) {
	// TODO: check if sub redirect and drop on >= 3

	trackingLink := r.trackingLinksRepository.FindTrackingLink(ctx, slug)
	if trackingLink == nil {
		return nil, TrackingLinkNotFoundError
	}

	if !trackingLink.IsActive {
		return nil, TrackingLinkDisabledError
	}

	if len(trackingLink.AllowedProtocols) > 0 && !contains(requestData.Protocol, trackingLink.AllowedProtocols) {
		return nil, UnsupportedProtocolError
	}

	countryCode, err := r.ipAddressParser.Parse(requestData.IP)
	if err != nil {
		r.log.Errorf("an error occured while parsing ip address (%s). error: %s\n", requestData.IP, err)
	}
	if len(trackingLink.AllowedGeos) > 0 && !contains(countryCode, trackingLink.AllowedGeos) {
		//TODO: handle trackingLink.UnsupportedGeoRedirectRules
		return nil, UnsupportedGeoError
	}

	ua, err := r.userAgentParser.Parse(requestData.UserAgent)
	if err != nil {
		r.log.Errorf("an error occured while parsing user-agent header (%s). error: %s\n", requestData.UserAgent, err)
	}
	if len(trackingLink.AllowedDevices) > 0 && !contains(ua.Device, trackingLink.AllowedDevices) {
		//TODO: handle trackingLink.UnsupportedDeviceRedirectRules
		return nil, UnsupportedDeviceError
	}
	//TODO: check OS + handle trackingLink.UnsupportedOSRedirectRules

	if trackingLink.IsCampaignOveraged {
		return r.handleRedirectRules(trackingLink.CampaignOverageRedirectRules, ctx, requestData, trackingLink, ua, countryCode)
	}

	if !trackingLink.IsCampaignActive {
		return r.handleRedirectRules(trackingLink.CampaignDisabledRedirectRules, ctx, requestData, trackingLink, ua, countryCode)
	}

	//TODO: prepare target URL template!
	//if deeplinkURL, ok := requestData.Params["deeplink"]; ok && trackingLink.AllowDeeplink {
	//	//TODO: handle deeplink
	//}

	targetURL := r.renderTokens(trackingLink, requestData, ua, countryCode)

	outputCh := r.registerClick(ctx, targetURL, trackingLink, requestData, ua, countryCode)

	return &dto.RedirectResult{
		TargetURL: targetURL,
		OutputCh:  outputCh,
	}, nil
}

func (r *redirectInteractor) handleRedirectRules(
	rr *valueobject.RedirectRules,
	ctx context.Context,
	requestData *dto.RedirectRequestData,
	trackingLink *entity.TrackingLink,
	userAgent *valueobject.UserAgent,
	countryCode string,
) (*dto.RedirectResult, error) {
	switch rr.RedirectType {
	case valueobject.LinkRedirectType:
		return &dto.RedirectResult{
			TargetURL: rr.RedirectURL,
			OutputCh:  r.registerClick(ctx, rr.RedirectURL, trackingLink, requestData, userAgent, countryCode),
		}, nil
	case valueobject.SlugRedirectType:
		return r.Redirect(ctx, rr.RedirectSlug, requestData)
	case valueobject.SmartSlugRedirectType:
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		newSlug := rr.RedirectSmartSlug[rnd.Intn(len(rr.RedirectSmartSlug))]
		return r.Redirect(ctx, newSlug, requestData)
	case valueobject.NoRedirectType:
		return nil, BlockRedirectError
	default:
		return nil, InvalidRedirectTypeError
	}
}

func (r *redirectInteractor) renderTokens(trackingLink *entity.TrackingLink, requestData *dto.RedirectRequestData, ua *valueobject.UserAgent, countryCode string) string {
	targetURL := trackingLink.TargetURLTemplate

	tokens := r.tokenRegExp.FindAllString(targetURL, -1)
	for _, token := range tokens {
		switch token {
		case IPAddressToken:
			targetURL = strings.ReplaceAll(targetURL, token, requestData.IP.String())
		case ClickIDToken:
			targetURL = strings.ReplaceAll(targetURL, token, requestData.RequestID)
		case UserAgentToken:
			targetURL = strings.ReplaceAll(targetURL, token, requestData.UserAgent)
		case CampaignIDToken:
			targetURL = strings.ReplaceAll(targetURL, token, trackingLink.CampaignID)
		case AffiliateIDToken:
			targetURL = strings.ReplaceAll(targetURL, token, trackingLink.AffiliateID)
		case SourceIDToken:
			targetURL = strings.ReplaceAll(targetURL, token, trackingLink.SourceID)
		case AdvertiserIDToken:
			targetURL = strings.ReplaceAll(targetURL, token, trackingLink.AdvertiserID)
		case DateToken:
			targetURL = strings.ReplaceAll(targetURL, token, time.Now().Format("2006-01-02"))
		case DateTimeToken:
			targetURL = strings.ReplaceAll(targetURL, token, time.Now().Format("2006-01-02T15:04:05"))
		case TimestampToken:
			targetURL = strings.ReplaceAll(targetURL, token, strconv.FormatInt(time.Now().Unix(), 10))
		case P1Token:
			values := requestData.GetParam("p1")
			targetURL = strings.ReplaceAll(targetURL, token, strings.Join(values, ","))
		case P2Token:
			values := requestData.GetParam("p2")
			targetURL = strings.ReplaceAll(targetURL, token, strings.Join(values, ","))
		case P3Token:
			values := requestData.GetParam("p3")
			targetURL = strings.ReplaceAll(targetURL, token, strings.Join(values, ","))
		case P4Token:
			values := requestData.GetParam("p4")
			targetURL = strings.ReplaceAll(targetURL, token, strings.Join(values, ","))
		case CountryCodeToken:
			targetURL = strings.ReplaceAll(targetURL, token, countryCode)
		case RefererToken:
			targetURL = strings.ReplaceAll(targetURL, token, requestData.Referer)
		case RandomStrToken:
			targetURL = strings.ReplaceAll(targetURL, token, randString(32))
		case RandomIntToken:
			targetURL = strings.ReplaceAll(targetURL, token, strconv.Itoa(rand.Intn(99999999-10000)+10000))
		case DeviceToken:
			targetURL = strings.ReplaceAll(targetURL, token, ua.Device)
		case PlatformToken:
			targetURL = strings.ReplaceAll(targetURL, token, ua.Platform)
		//TODO: add new tokens

		// replace undefined tokens with empty string
		default:
			targetURL = strings.ReplaceAll(targetURL, token, "")
		}
	}

	return targetURL
}

func (r *redirectInteractor) registerClick(ctx context.Context, targetURL string, trackingLink *entity.TrackingLink, requestData *dto.RedirectRequestData, ua *valueobject.UserAgent, countryCode string) <-chan *dto.ClickProcessingResult {
	click := &entity.Click{
		ID:          requestData.RequestID,
		TargetURL:   targetURL,
		TRKLink:     trackingLink,
		UserAgent:   ua,
		CountryCode: countryCode,
		IP:          requestData.IP,
		Referer:     requestData.Referer,

		P1: strings.Join(requestData.GetParam("p1"), ","),
		P2: strings.Join(requestData.GetParam("p2"), ","),
		P3: strings.Join(requestData.GetParam("p3"), ","),
		P4: strings.Join(requestData.GetParam("p4"), ","),
	}

	outputs := make([]<-chan *dto.ClickProcessingResult, len(r.clickHandlers))
	for _, handler := range r.clickHandlers {
		outputs = append(outputs, handler.HandleClick(ctx, click))
	}

	return merge(outputs...)
}

// merge function will fan-in the results received from ClickHandlerInterface(s).
func merge(clkProcessingResultChans ...<-chan *dto.ClickProcessingResult) <-chan *dto.ClickProcessingResult {
	var wg sync.WaitGroup
	out := make(chan *dto.ClickProcessingResult)

	mergeFunc := func(c <-chan *dto.ClickProcessingResult) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(clkProcessingResultChans))
	for _, c := range clkProcessingResultChans {
		go mergeFunc(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
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
