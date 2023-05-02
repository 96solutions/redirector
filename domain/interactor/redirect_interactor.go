package interactor

import (
	"context"
	"errors"
	"log"
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

type redirectResult struct {
	TargetURL string
	OutputCh  <-chan *ClickProcessingResult
}

type ClickProcessingResult struct {
	Click *entity.Click
	Err   error
}

//go:generate mockgen -package=mocks -destination=mocks/mock_redirect_interactor.go -source=domain/interactor/redirect_interactor.go RedirectInteractor

// RedirectInteractor interface describes the service to handle requests and return the following target URL to redirect to.
type RedirectInteractor interface {
	Redirect(ctx context.Context, slug string, requestData *dto.RedirectRequestData) (*redirectResult, error)
}

type redirectInteractor struct {
	trackingLinksRepository repository.TrackingLinksRepositoryInterface
	ipAddressParser         service.IpAddressParserInterface
	userAgentParser         service.UserAgentParser
	tokenRegExp             *regexp.Regexp
	clickHandlers           []ClickHandlerInterface
}

// NewRedirectInteractor function creates RedirectInteractor implementation.
func NewRedirectInteractor(
	trkRepo repository.TrackingLinksRepositoryInterface,
	ipAddressParser service.IpAddressParserInterface,
	userAgentParser service.UserAgentParser,
	clickHandlers []ClickHandlerInterface,
) RedirectInteractor {
	compiledRegExp, err := regexp.Compile(`{({)?(\w+)(})?}`)
	if err != nil {
		panic(err)
	}

	return &redirectInteractor{
		trackingLinksRepository: trkRepo,
		ipAddressParser:         ipAddressParser,
		userAgentParser:         userAgentParser,
		tokenRegExp:             compiledRegExp,
		clickHandlers:           clickHandlers,
	}
}

// Redirect function handles requests and returns the target URL to redirect traffic to.
func (r *redirectInteractor) Redirect(ctx context.Context, slug string, requestData *dto.RedirectRequestData) (*redirectResult, error) {
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
		log.Printf("an error occured while parsing ip address (%s). error: %s\n", requestData.IP, err)
	}
	if len(trackingLink.AllowedGeos) > 0 && !contains(countryCode, trackingLink.AllowedGeos) {
		return nil, UnsupportedGeoError
	}

	ua, err := r.userAgentParser.Parse(requestData.UserAgent)
	if err != nil {
		log.Printf("an error occured while parsing user-agent header (%s). error: %s\n", requestData.UserAgent, err)
	}
	if len(trackingLink.AllowedDevices) > 0 && !contains(ua.Device, trackingLink.AllowedDevices) {
		return nil, UnsupportedDeviceError
	}

	if trackingLink.IsCampaignOveraged {
		return r.handleRedirectRules(trackingLink.CampaignOverageRedirectRules, ctx, requestData, trackingLink, ua, countryCode)
	}

	if !trackingLink.IsCampaignActive {
		return r.handleRedirectRules(trackingLink.CampaignDisabledRedirectRules, ctx, requestData, trackingLink, ua, countryCode)
	}

	//if deeplinkURL, ok := requestData.Params["deeplink"]; ok && trackingLink.AllowDeeplink {
	//	//TODO: handle deeplink
	//}

	targetURL := r.renderTokens(trackingLink, requestData, ua, countryCode)

	outputCh := r.registerClick(ctx, targetURL, trackingLink, requestData, ua, countryCode)

	return &redirectResult{
		TargetURL: targetURL,
		OutputCh:  outputCh,
	}, nil
}

func (r *redirectInteractor) handleRedirectRules(rr *valueobject.RedirectRules, ctx context.Context, requestData *dto.RedirectRequestData, trackingLink *entity.TrackingLink, userAgent *valueobject.UserAgent, countryCode string) (*redirectResult, error) {
	switch rr.RedirectType {
	case valueobject.LinkRedirectType:
		return &redirectResult{
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

func (r *redirectInteractor) registerClick(ctx context.Context, targetURL string, trackingLink *entity.TrackingLink, requestData *dto.RedirectRequestData, ua *valueobject.UserAgent, countryCode string) <-chan *ClickProcessingResult {
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

	outputs := make([]<-chan *ClickProcessingResult, len(r.clickHandlers))
	for _, handler := range r.clickHandlers {
		outputs = append(outputs, handler.HandleClick(ctx, click))
	}

	return merge(outputs...)
}

// merge function will fan-in the results received from ClickHandlerInterface(s).
func merge(cs ...<-chan *ClickProcessingResult) <-chan *ClickProcessingResult {
	var wg sync.WaitGroup
	out := make(chan *ClickProcessingResult)

	output := func(c <-chan *ClickProcessingResult) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
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
