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
	// ErrUnsupportedProtocol is returned when the request protocol is not allowed.
	ErrUnsupportedProtocol = errors.New("protocol is not allowed for that tracking link")
	// ErrUnsupportedGeo is returned when the visitor's geo location is not allowed.
	ErrUnsupportedGeo = errors.New("visitor geo is not allowed for that tracking link")
	// ErrUnsupportedDevice is returned when the visitor's device type is not allowed.
	ErrUnsupportedDevice = errors.New("visitor device is not allowed for that tracking link")
	// ErrUnsupportedOS is returned when the visitor's operating system is not allowed.
	ErrUnsupportedOS = errors.New("visitor OS is not allowed for that tracking link")
	// ErrInvalidRedirectType is returned when the redirect rules contain an invalid type.
	ErrInvalidRedirectType = errors.New("invalid redirect type is stored in tracking link redirect rules")
	// ErrBlockRedirect is returned when the redirect should be blocked.
	ErrBlockRedirect = errors.New("redirect should be blocked")
	// ErrTrackingLinkDisabled is returned when the tracking link is disabled.
	ErrTrackingLinkDisabled = errors.New("used tracking link is disabled")
	// ErrTrackingLinkNotFound is returned when no tracking link is found for the given slug.
	ErrTrackingLinkNotFound = errors.New("no tracking link was found by slug")
)

const (
	ipAddressToken    = "{ip}"
	clickIDToken      = "{click_id}"
	userAgentToken    = "{user_agent}"
	campaignIDToken   = "{campaign_id}"
	affiliateIDToken  = "{aff_id}"
	sourceIDToken     = "{source_id}"
	advertiserIDToken = "{advertiser_id}"
	dateToken         = "{date}"
	dateTimeToken     = "{date_time}"
	timestampToken    = "{timestamp}"
	p1Token           = "{p1}"
	p2Token           = "{p2}"
	p3Token           = "{p3}"
	p4Token           = "{p4}"
	countryCodeToken  = "{country_code}"
	refererToken      = "{referer}"
	randomStrToken    = "{random_str}"
	randomIntToken    = "{random_int}"
	deviceToken       = "{device}"
	platformToken     = "{platform}"

	unknownStrValue = "unknown"

	randomStringLen = 32
)

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
// including tracking link validation, click tracking, and redirect rule application.
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
	trackingLink := r.trackingLinksRepository.FindTrackingLink(ctx, slug)
	if trackingLink == nil {
		return nil, ErrTrackingLinkNotFound
	}

	if !trackingLink.IsActive {
		return nil, ErrTrackingLinkDisabled
	}

	if len(trackingLink.AllowedProtocols) > 0 && !trackingLink.AllowedProtocols[requestData.Protocol] {
		return nil, ErrUnsupportedProtocol
	}

	countryCode, err := r.ipAddressParser.Parse(requestData.IP)
	if err != nil {
		slog.Error("an error occurred while parsing ip address", "ip", requestData.IP, "error", err)
		countryCode = unknownStrValue
	}

	ua, err := r.userAgentParser.Parse(requestData.UserAgent)
	if err != nil {
		slog.Error("an error occurred while parsing user-agent header", "user-agent", requestData.UserAgent, "error", err)
		ua = &valueobject.UserAgent{
			SrcString: requestData.UserAgent,
			Device:    unknownStrValue,
			Platform:  unknownStrValue,
			Browser:   unknownStrValue,
		}
	}

	if len(trackingLink.AllowedGeos) > 0 && !trackingLink.AllowedGeos[countryCode] {
		return r.handleRedirectRules(
			ctx,
			trackingLink.CampaignOverageRedirectRules,
			requestData,
			trackingLink,
			countryCode,
			ua,
			ErrUnsupportedGeo,
		)
	}

	if len(trackingLink.AllowedDevices) > 0 && !trackingLink.AllowedDevices[ua.Device] {
		return r.handleRedirectRules(
			ctx,
			trackingLink.CampaignOverageRedirectRules,
			requestData,
			trackingLink,
			countryCode,
			ua,
			ErrUnsupportedDevice,
		)
	}

	if len(trackingLink.AllowedOS) > 0 && !trackingLink.AllowedOS[ua.Platform] {
		return r.handleRedirectRules(
			ctx,
			trackingLink.CampaignOverageRedirectRules,
			requestData,
			trackingLink,
			countryCode,
			ua,
			ErrUnsupportedOS,
		)
	}

	if trackingLink.IsCampaignOveraged {
		return r.handleRedirectRules(
			ctx,
			trackingLink.CampaignOverageRedirectRules,
			requestData,
			trackingLink,
			countryCode,
			ua,
			nil,
		)
	}

	if !trackingLink.IsCampaignActive {
		return r.handleRedirectRules(
			ctx,
			trackingLink.CampaignDisabledRedirectRules,
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
	ctx context.Context,
	rr *valueobject.RedirectRules,
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

		return nil, ErrBlockRedirect
	default:
		return nil, ErrInvalidRedirectType
	}
}

func (r *redirectInteractor) makeRedirectTemplate(
	trackingLink *entity.TrackingLink,
	requestData *dto.RedirectRequestData,
) string {
	targetURL := trackingLink.TargetURLTemplate

	if landingURL, paramExists := requestData.Params["landing"]; paramExists {
		if landing, landingExists := trackingLink.LandingPages[landingURL[0]]; landingExists {
			targetURL = landing.TargetURL
		}
	}

	if deeplinkURL, ok := requestData.Params["deeplink"]; ok && trackingLink.AllowDeeplink {
		targetURL = deeplinkURL[0]
	}

	return targetURL
}

func (r *redirectInteractor) renderTokens(trackingLink *entity.TrackingLink, requestData *dto.RedirectRequestData, ua *valueobject.UserAgent, countryCode string) string {
	targetURL := r.makeRedirectTemplate(trackingLink, requestData)

	tokens := r.tokenRegExp.FindAllString(targetURL, -1)
	for _, token := range tokens {
		switch token {
		case ipAddressToken:
			targetURL = strings.ReplaceAll(targetURL, token, requestData.IP.String())
		case clickIDToken:
			targetURL = strings.ReplaceAll(targetURL, token, requestData.RequestID)
		case userAgentToken:
			targetURL = strings.ReplaceAll(targetURL, token, requestData.UserAgent)
		case campaignIDToken:
			targetURL = strings.ReplaceAll(targetURL, token, trackingLink.CampaignID)
		case affiliateIDToken:
			targetURL = strings.ReplaceAll(targetURL, token, trackingLink.AffiliateID)
		case sourceIDToken:
			targetURL = strings.ReplaceAll(targetURL, token, trackingLink.SourceID)
		case advertiserIDToken:
			targetURL = strings.ReplaceAll(targetURL, token, trackingLink.AdvertiserID)
		case dateToken:
			targetURL = strings.ReplaceAll(targetURL, token, time.Now().Format("2006-01-02"))
		case dateTimeToken:
			targetURL = strings.ReplaceAll(targetURL, token, time.Now().Format("2006-01-02T15:04:05"))
		case timestampToken:
			targetURL = strings.ReplaceAll(targetURL, token, strconv.FormatInt(time.Now().Unix(), 10))
		case p1Token:
			values := requestData.GetParam("p1")
			targetURL = strings.ReplaceAll(targetURL, token, strings.Join(values, ","))
		case p2Token:
			values := requestData.GetParam("p2")
			targetURL = strings.ReplaceAll(targetURL, token, strings.Join(values, ","))
		case p3Token:
			values := requestData.GetParam("p3")
			targetURL = strings.ReplaceAll(targetURL, token, strings.Join(values, ","))
		case p4Token:
			values := requestData.GetParam("p4")
			targetURL = strings.ReplaceAll(targetURL, token, strings.Join(values, ","))
		case countryCodeToken:
			targetURL = strings.ReplaceAll(targetURL, token, countryCode)
		case refererToken:
			targetURL = strings.ReplaceAll(targetURL, token, requestData.Referer)
		case randomStrToken:
			targetURL = strings.ReplaceAll(targetURL, token, randString(randomStringLen))
		case randomIntToken:
			targetURL = strings.ReplaceAll(targetURL, token, strconv.Itoa(rand.Intn(99999999-10000)+10000))
		case deviceToken:
			targetURL = strings.ReplaceAll(targetURL, token, ua.Device)
		case platformToken:
			targetURL = strings.ReplaceAll(targetURL, token, ua.Platform)

		// replace undefined tokens with empty string
		default:
			targetURL = strings.ReplaceAll(targetURL, token, "")
		}
	}

	//TODO: append gclid query param if present in requestData.Params

	return targetURL
}

func (r *redirectInteractor) registerClick(
	ctx context.Context,
	slug string,
	targetURL string,
	trackingLink *entity.TrackingLink,
	requestData *dto.RedirectRequestData,
	ua *valueobject.UserAgent,
	countryCode string,
) <-chan *dto.ClickProcessingResult {
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randString function generates random string of n-length.
func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
