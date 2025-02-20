// Package interactor contains all use-case interactors preformed by the application.
package interactor

import (
	"context"
	"errors"
	"log/slog"
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
	UnsupportedOSError        = errors.New("visitor OS is not allowed for that tracking link")
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

	UnknownStrValue = "unknown"
)

//mockgen -package=mocks -destination=mocks/mock_redirect_interactor.go -source=domain/interactor/redirect_interactor.go RedirectInteractor
//go:generate mockgen -package=mocks -destination=mocks/mock_redirect_interactor.go -source=redirect_interactor.go RedirectInteractor

// RedirectInteractor handles the business logic for processing redirect requests.
// It validates tracking links, applies redirect rules, and manages click tracking.
type RedirectInteractor interface {
	// Redirect processes a redirect request for a given slug and request data.
	// It validates the tracking link, checks various conditions (protocol, geo, device),
	// applies redirect rules, and returns the target URL along with click tracking results.
	// Returns an error if the redirect cannot be processed.
	Redirect(ctx context.Context, slug string, requestData *dto.RedirectRequestData) (*dto.RedirectResult, error)
}

// redirectInteractor implements RedirectInteractor interface and handles the core redirect logic
// including tracking link validation, click tracking, and redirect rule application
type redirectInteractor struct {
	trackingLinksRepository repository.TrackingLinksRepositoryInterface
	ipAddressParser         service.IPAddressParserInterface
	userAgentParser         service.UserAgentParserInterface
	tokenRegExp             *regexp.Regexp
	clickHandlers           []ClickHandlerInterface
}

// NewRedirectInteractor function creates RedirectInteractor implementation.
func NewRedirectInteractor(
	trkRepo repository.TrackingLinksRepositoryInterface,
	ipAddressParser service.IPAddressParserInterface,
	userAgentParser service.UserAgentParserInterface,
	clickHandlers []ClickHandlerInterface,
) RedirectInteractor {
	compiledRegExp := regexp.MustCompile(`{({)?(\w+)(})?}`)

	return &redirectInteractor{
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

	if len(trackingLink.AllowedProtocols) > 0 && !trackingLink.AllowedProtocols[requestData.Protocol] {
		return nil, UnsupportedProtocolError
	}

	countryCode, err := r.ipAddressParser.Parse(requestData.IP)
	if err != nil {
		slog.Error("an error occurred while parsing ip address", "ip", requestData.IP, "error", err)
		countryCode = UnknownStrValue
	}

	ua, err := r.userAgentParser.Parse(requestData.UserAgent)
	if err != nil {
		slog.Error("an error occurred while parsing user-agent header", "user-agent", requestData.UserAgent, "error", err)
		ua = &valueobject.UserAgent{
			SrcString: requestData.UserAgent,
			Device:    UnknownStrValue,
			Platform:  UnknownStrValue,
			Browser:   UnknownStrValue,
		}
	}

	if len(trackingLink.AllowedGeos) > 0 && !trackingLink.AllowedGeos[countryCode] {
		return r.handleRedirectRules(
			trackingLink.CampaignOverageRedirectRules,
			ctx,
			requestData,
			trackingLink,
			countryCode,
			ua,
			UnsupportedGeoError,
		)
	}

	if len(trackingLink.AllowedDevices) > 0 && !trackingLink.AllowedDevices[ua.Device] {
		return r.handleRedirectRules(
			trackingLink.CampaignOverageRedirectRules,
			ctx,
			requestData,
			trackingLink,
			countryCode,
			ua,
			UnsupportedDeviceError,
		)
	}

	if len(trackingLink.AllowedOS) > 0 && !trackingLink.AllowedOS[ua.Platform] {
		return r.handleRedirectRules(
			trackingLink.CampaignOverageRedirectRules,
			ctx,
			requestData,
			trackingLink,
			countryCode,
			ua,
			UnsupportedOSError,
		)
	}

	if trackingLink.IsCampaignOveraged {
		return r.handleRedirectRules(
			trackingLink.CampaignOverageRedirectRules,
			ctx,
			requestData,
			trackingLink,
			countryCode,
			ua,
			nil,
		)
	}

	if !trackingLink.IsCampaignActive {
		return r.handleRedirectRules(
			trackingLink.CampaignDisabledRedirectRules,
			ctx,
			requestData,
			trackingLink,
			countryCode,
			ua,
			nil,
		)
	}

	targetURL := r.renderTokens(trackingLink, requestData, ua, countryCode)

	outputCh := r.registerClick(ctx, slug, targetURL, trackingLink, requestData, ua, countryCode)

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
	countryCode string,
	userAgent *valueobject.UserAgent,
	err error,
) (*dto.RedirectResult, error) {
	switch rr.RedirectType {
	case valueobject.LinkRedirectType:
		return &dto.RedirectResult{
			TargetURL: rr.RedirectURL,
			OutputCh:  r.registerClick(ctx, requestData.Slug, rr.RedirectURL, trackingLink, requestData, userAgent, countryCode),
		}, nil
	case valueobject.SlugRedirectType:
		ctx = context.WithValue(ctx, "slug", requestData.Slug)
		return r.Redirect(ctx, rr.RedirectSlug, requestData)
	case valueobject.SmartSlugRedirectType:
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		newSlug := rr.RedirectSmartSlug[rnd.Intn(len(rr.RedirectSmartSlug))]
		ctx = context.WithValue(ctx, "slug", requestData.Slug)
		return r.Redirect(ctx, newSlug, requestData)
	case valueobject.NoRedirectType:
		if err != nil {
			return nil, err
		}

		return nil, BlockRedirectError
	default:
		return nil, InvalidRedirectTypeError
	}
}

func (r *redirectInteractor) renderTokens(trackingLink *entity.TrackingLink, requestData *dto.RedirectRequestData, ua *valueobject.UserAgent, countryCode string) string {
	targetURL := trackingLink.TargetURLTemplate

	if landingURL, paramExists := requestData.Params["landing"]; paramExists {
		if landing, landingExists := trackingLink.LandingPages[landingURL[0]]; landingExists {
			targetURL = landing.TargetURL
		}
	}

	if deeplinkURL, ok := requestData.Params["deeplink"]; ok && trackingLink.AllowDeeplink {
		targetURL = deeplinkURL[0]
	}

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

		// replace undefined tokens with empty string
		default:
			targetURL = strings.ReplaceAll(targetURL, token, "")
		}
	}

	//TODO: append gclid query param if present in requestData.Params

	return targetURL
}

func (r *redirectInteractor) registerClick(ctx context.Context, slug string, targetURL string, trackingLink *entity.TrackingLink, requestData *dto.RedirectRequestData, ua *valueobject.UserAgent, countryCode string) <-chan *dto.ClickProcessingResult {
	click := &entity.Click{
		ID:           requestData.RequestID,
		TargetURL:    targetURL,
		Referer:      requestData.Referer,
		TrkURL:       requestData.URL.String(),
		ParentSlug:   ctx.Value("slug").(string),
		Slug:         slug,
		TRKLink:      trackingLink,
		SourceID:     trackingLink.SourceID,
		CampaignID:   trackingLink.CampaignID,
		AffiliateID:  trackingLink.AffiliateID,
		AdvertiserID: trackingLink.AdvertiserID,
		IsParallel:   false,
		LandingID:    requestData.Params["landing"][0],
		GCLID:        requestData.Params["gclid"][0],
		UserAgent:    ua,
		Agent:        ua.SrcString,
		Platform:     ua.Platform,
		Browser:      ua.Browser,
		Device:       ua.Device,
		IP:           requestData.IP,
		//Region:       "",
		CountryCode: countryCode,
		//City:         "",
		P1:        strings.Join(requestData.GetParam("p1"), ","),
		P2:        strings.Join(requestData.GetParam("p2"), ","),
		P3:        strings.Join(requestData.GetParam("p3"), ","),
		P4:        strings.Join(requestData.GetParam("p4"), ","),
		CreatedAt: time.Now(),
	}

	outputs := make([]<-chan *dto.ClickProcessingResult, len(r.clickHandlers))
	for i, handler := range r.clickHandlers {
		outputs[i] = handler.HandleClick(ctx, click)
	}

	return merge(outputs)
}

// merge function will fan-in the results received from ClickHandlerInterface(s).
func merge(clkProcessingResultChans []<-chan *dto.ClickProcessingResult) <-chan *dto.ClickProcessingResult {
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

// randString function generates random string of n-length.
func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
