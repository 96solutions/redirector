package interactor_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lroman242/redirector/domain/dto"
	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/domain/interactor"
	"github.com/lroman242/redirector/domain/valueobject"
	"github.com/lroman242/redirector/mocks"
	"go.uber.org/mock/gomock"
)

const (
	requestSlug = "test-slug"
	redirectURL = "https://example.com"
	ipAddress   = "192.168.1.1"
	userAgent   = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36"
	countryCode = "US"
	protocol    = "http"
	device      = "Mobile"
	platform    = "Android"
	browser     = "Chrome"
	requestID   = "1234567890"
	p1          = "param1"
	p2          = "param2"
	p3          = "param3"
	p4          = "param4"
	referrer    = "https://referrer.com"
)

// testData holds common test data to reduce duplication
type testData struct {
	slug        string
	requestData *dto.RedirectRequestData
	userAgent   *valueobject.UserAgent
	countryCode string
}

// setupTest creates a new test environment with mocked dependencies
func setupTest(t *testing.T) (
	*gomock.Controller,
	interactor.RedirectInteractor,
	*mocks.MockTrackingLinksRepositoryInterface,
	*mocks.MockIPAddressParserInterface,
	*mocks.MockUserAgentParserInterface,
	*mocks.MockClicksRepository,
) {
	ctrl := gomock.NewController(t)
	clkRepo := mocks.NewMockClicksRepository(ctrl)
	srv, trkRepo, ipParser, uaParser := makeRedirectInteractor(ctrl, interactor.NewStoreClickHandler(clkRepo))
	return ctrl, srv, trkRepo, ipParser, uaParser, clkRepo
}

// newTestData creates common test data
func newTestData() *testData {
	incomeURL, _ := url.Parse(redirectURL + "/" + requestSlug)
	return &testData{
		slug: requestSlug,
		requestData: &dto.RedirectRequestData{
			RequestID: requestID,
			Slug:      requestSlug,
			Params:    map[string][]string{"p1": []string{p1}, "p2": []string{p2}, "p3": []string{p3}, "p4": []string{p4}},
			Headers:   make(map[string][]string),
			UserAgent: userAgent,
			IP:        net.ParseIP(ipAddress),
			Protocol:  protocol,
			URL:       incomeURL,
			Referer:   referrer,
		},
		userAgent: &valueobject.UserAgent{
			Bot:      false,
			Device:   device,
			Platform: platform,
			Browser:  browser,
		},
		countryCode: countryCode,
	}
}

// makeRedirectInteractor creates a new RedirectInteractor with mocked dependencies
func makeRedirectInteractor(ctrl *gomock.Controller, handlers ...interactor.ClickHandlerInterface) (
	interactor.RedirectInteractor,
	*mocks.MockTrackingLinksRepositoryInterface,
	*mocks.MockIPAddressParserInterface,
	*mocks.MockUserAgentParserInterface,
) {
	trkRepo := mocks.NewMockTrackingLinksRepositoryInterface(ctrl)
	ipParser := mocks.NewMockIPAddressParserInterface(ctrl)
	uaParser := mocks.NewMockUserAgentParserInterface(ctrl)

	srv := interactor.NewRedirectInteractor(trkRepo, ipParser, uaParser, handlers)
	return srv, trkRepo, ipParser, uaParser
}

func TestRedirectInteractor_Redirect_TrackingLinkNotFound(t *testing.T) {
	ctrl, srv, trkRepo, _, _, _ := setupTest(t)
	defer ctrl.Finish()

	td := newTestData()
	trkRepo.EXPECT().FindTrackingLink(context.Background(), td.slug).Return(nil)

	result, err := srv.Redirect(context.Background(), td.slug, td.requestData)

	if !errors.Is(err, interactor.ErrTrackingLinkNotFound) {
		t.Error("expected TrackingLinkNotFound error")
	}
	if result != nil {
		t.Error("expected nil result")
	}
}

func TestRedirectInteractor_Redirect_DisabledTrackingLink(t *testing.T) {
	ctrl, srv, trkRepo, _, _, _ := setupTest(t)
	defer ctrl.Finish()

	td := newTestData()
	trkLink := &entity.TrackingLink{
		IsActive:         false,
		Slug:             td.slug,
		AllowedProtocols: entity.AllowedListType{"https": true},
	}

	trkRepo.EXPECT().FindTrackingLink(context.Background(), td.slug).Return(trkLink)

	result, err := srv.Redirect(context.Background(), td.slug, td.requestData)

	if !errors.Is(err, interactor.ErrTrackingLinkDisabled) {
		t.Error("expected TrackingLinkDisabled error")
	}
	if result != nil {
		t.Error("expected nil result")
	}
}

func TestRedirectInteractor_Redirect_WrongProtocolError(t *testing.T) {
	ctrl, srv, trkRepo, _, _, _ := setupTest(t)
	defer ctrl.Finish()

	td := newTestData()
	td.requestData.Protocol = "http"

	trkLink := &entity.TrackingLink{
		IsActive:         true,
		Slug:             td.slug,
		AllowedProtocols: entity.AllowedListType{"https": true},
	}

	trkRepo.EXPECT().FindTrackingLink(context.Background(), td.slug).Return(trkLink)

	result, err := srv.Redirect(context.Background(), td.slug, td.requestData)

	if !errors.Is(err, interactor.ErrUnsupportedProtocol) {
		t.Error("expected UnsupportedProtocol error")
	}
	if result != nil {
		t.Error("expected nil result")
	}
}

func TestRedirectInteractor_Redirect_WrongGeoError(t *testing.T) {
	ctrl, srv, trkRepo, ipParser, uaParser, _ := setupTest(t)
	defer ctrl.Finish()

	td := newTestData()

	trkLink := &entity.TrackingLink{
		IsActive:         true,
		Slug:             td.slug,
		AllowedProtocols: make(entity.AllowedListType),
		AllowedGeos:      entity.AllowedListType{"US": true, "PT": true, "UA": true},
		CampaignGeoRedirectRules: &valueobject.RedirectRules{
			RedirectType: valueobject.NoRedirectType,
		},
	}

	trkRepo.EXPECT().FindTrackingLink(context.Background(), td.slug).Return(trkLink)
	ipParser.EXPECT().Parse(td.requestData.IP).Return("PL", nil)
	uaParser.EXPECT().Parse(td.requestData.UserAgent).Return(td.userAgent, nil)

	result, err := srv.Redirect(context.Background(), td.slug, td.requestData)

	if !errors.Is(err, interactor.ErrUnsupportedGeo) {
		t.Error("expected UnsupportedGeo error")
	}
	if result != nil {
		t.Error("expected nil result")
	}
}

func TestRedirectInteractor_Redirect_WrongDeviceError(t *testing.T) {
	ctrl, srv, trkRepo, ipParser, uaParser, _ := setupTest(t)
	defer ctrl.Finish()

	td := newTestData()

	trkLink := &entity.TrackingLink{
		IsActive:         true,
		Slug:             td.slug,
		AllowedProtocols: make(entity.AllowedListType),
		AllowedGeos:      make(entity.AllowedListType),
		AllowedDevices:   entity.AllowedListType{"desktop": true},
		CampaignDevicesRedirectRules: &valueobject.RedirectRules{
			RedirectType: valueobject.NoRedirectType,
		},
	}

	trkRepo.EXPECT().FindTrackingLink(context.Background(), td.slug).Return(trkLink)
	ipParser.EXPECT().Parse(td.requestData.IP).Return(td.countryCode, nil)
	uaParser.EXPECT().Parse(td.requestData.UserAgent).Return(td.userAgent, nil)

	result, err := srv.Redirect(context.Background(), td.slug, td.requestData)

	if !errors.Is(err, interactor.ErrUnsupportedDevice) {
		t.Error("expected UnsupportedDevice error")
	}
	if result != nil {
		t.Error("expected nil result")
	}
}

func TestRedirectInteractor_Redirect_CampaignOveraged(t *testing.T) {
	testCases := []struct {
		name              string
		trkLink           *entity.TrackingLink
		expectedTargetURL string
		expectedError     error
	}{
		{
			name: "OveragedRedirectRulesLinkRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				IsCampaignActive:   true,
				AllowedProtocols:   entity.AllowedListType{},
				AllowedGeos:        entity.AllowedListType{},
				AllowedDevices:     entity.AllowedListType{},
				IsCampaignOveraged: true,
				CampaignOverageRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.LinkRedirectType,
					RedirectURL:       "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignOveraged/OveragedRedirectRulesLinkRedirectType",
					RedirectSlug:      "",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignOveraged/OveragedRedirectRulesLinkRedirectType",
			expectedError:     nil,
		},
		{
			name: "OveragedRedirectRulesSlugRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				IsCampaignActive:   true,
				AllowedProtocols:   entity.AllowedListType{},
				AllowedGeos:        entity.AllowedListType{},
				AllowedDevices:     entity.AllowedListType{},
				IsCampaignOveraged: true,
				CampaignOverageRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.SlugRedirectType,
					RedirectURL:       "",
					RedirectSlug:      "testSlug456",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignOveraged/OveragedRedirectRulesSlugRedirectType",
			expectedError:     nil,
		},
		{
			name: "OveragedNoRedirectRule",
			trkLink: &entity.TrackingLink{
				IsActive:                     true,
				IsCampaignActive:             true,
				AllowedProtocols:             entity.AllowedListType{},
				AllowedGeos:                  entity.AllowedListType{},
				AllowedDevices:               entity.AllowedListType{},
				IsCampaignOveraged:           true,
				CampaignOverageRedirectRules: nil,
			},
			expectedTargetURL: "",
			expectedError:     interactor.ErrInvalidRedirectRules,
		},
		{
			name: "OveragedRedirectRulesSmartSlugRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				IsCampaignActive:   true,
				AllowedProtocols:   entity.AllowedListType{},
				AllowedGeos:        entity.AllowedListType{},
				AllowedDevices:     entity.AllowedListType{},
				IsCampaignOveraged: true,
				CampaignOverageRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.SmartSlugRedirectType,
					RedirectURL:       "",
					RedirectSlug:      "",
					RedirectSmartSlug: []string{"testSlug000", "testSlug111", "testSlug222"},
				},
			},
			expectedTargetURL: "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignOveraged/OveragedRedirectRulesSmartSlugRedirectType",
			expectedError:     nil,
		},
		{
			name: "OveragedRedirectRulesNoRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				IsCampaignActive:   true,
				AllowedProtocols:   map[string]bool{},
				AllowedGeos:        map[string]bool{},
				AllowedDevices:     map[string]bool{},
				IsCampaignOveraged: true,
				CampaignOverageRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.NoRedirectType,
					RedirectURL:       "",
					RedirectSlug:      "",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "",
			expectedError:     interactor.ErrBlockRedirect,
		},
		{
			name: "OveragedRedirectRulesInvalidRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				IsCampaignActive:   true,
				AllowedProtocols:   map[string]bool{},
				AllowedGeos:        map[string]bool{},
				AllowedDevices:     map[string]bool{},
				IsCampaignOveraged: true,
				CampaignOverageRedirectRules: &valueobject.RedirectRules{
					RedirectType:      "invalid_redirect_type",
					RedirectURL:       "",
					RedirectSlug:      "",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "",
			expectedError:     interactor.ErrInvalidRedirectType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, srv, trkRepo, ipParser, uaParser, clkRepo := setupTest(t)
			defer ctrl.Finish()

			td := newTestData()
			tc.trkLink.Slug = td.slug

			trkRepo.EXPECT().FindTrackingLink(context.Background(), td.slug).Return(tc.trkLink)
			ipParser.EXPECT().Parse(td.requestData.IP).Return(td.countryCode, nil)
			uaParser.EXPECT().Parse(td.requestData.UserAgent).Return(td.userAgent, nil)

			if tc.expectedError == nil {
				clkRepo.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					// Return(nil)
					DoAndReturn(func(_ context.Context, click *entity.Click) error {
						if click.ID != td.requestData.RequestID {
							t.Errorf("unexpected click ID. expected %s but got %s", td.requestData.RequestID, click.ID)
						}
						if click.TargetURL != tc.expectedTargetURL {
							t.Errorf("unexpected target URL. expected %s but got %s", tc.expectedTargetURL, click.TargetURL)
						}

						return nil
					})
			}

			if tc.trkLink.CampaignOverageRedirectRules != nil &&
				(tc.trkLink.CampaignOverageRedirectRules.RedirectType == valueobject.SlugRedirectType ||
					tc.trkLink.CampaignOverageRedirectRules.RedirectType == valueobject.SmartSlugRedirectType) {
				trkRepo.
					EXPECT().
					FindTrackingLink(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_, arg interface{}) *entity.TrackingLink {
						slug, ok := arg.(string)
						if !ok {
							t.Error("invalid argument type. expected string")
						}

						if tc.trkLink.CampaignOverageRedirectRules.RedirectType == valueobject.SlugRedirectType {
							if slug != tc.trkLink.CampaignOverageRedirectRules.RedirectSlug {
								t.Errorf("invalid argument received. expected %s but got %s", tc.trkLink.CampaignOverageRedirectRules.RedirectSlug, slug)
							}
						} else if tc.trkLink.CampaignOverageRedirectRules.RedirectType == valueobject.SmartSlugRedirectType {
							inArray := false

							for _, sl := range tc.trkLink.CampaignOverageRedirectRules.RedirectSmartSlug {
								if sl == slug {
									inArray = true
								}
							}

							if !inArray {
								t.Errorf("invalid argument received. expected one of %v but got %s", tc.trkLink.CampaignOverageRedirectRules.RedirectSmartSlug, slug)
							}
						}

						return &entity.TrackingLink{
							IsActive:           true,
							Slug:               slug,
							AllowedProtocols:   map[string]bool{},
							AllowedGeos:        map[string]bool{},
							AllowedDevices:     map[string]bool{},
							IsCampaignOveraged: false,
							IsCampaignActive:   true,
							TargetURLTemplate:  "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignOveraged/" + tc.name,
						}
					})
				ipParser.EXPECT().Parse(td.requestData.IP).Return(td.countryCode, nil)
				uaParser.EXPECT().Parse(td.requestData.UserAgent).Return(td.userAgent, nil)
			}

			rResult, err := srv.Redirect(context.Background(), td.slug, td.requestData)
			if tc.expectedError != nil {
				if !errors.Is(err, tc.expectedError) {
					t.Errorf("unexpected result, %T expected", tc.expectedError)
				}
				if rResult != nil {
					t.Error("unexpected target url. expected empty value")
				}
			} else if tc.expectedTargetURL != rResult.TargetURL {
				t.Errorf("unexpected target url. expected %s but got %s", tc.expectedTargetURL, rResult.TargetURL)
			} else {
				<-rResult.OutputCh
			}
		})
	}
}

func TestRedirectInteractor_Redirect_CampaignDisabled(t *testing.T) {
	testCases := []struct {
		name              string
		trkLink           *entity.TrackingLink
		expectedTargetURL string
		expectedError     error
	}{
		{
			name: "CampaignDisableRedirectRulesLinkRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				AllowedProtocols:   entity.AllowedListType{},
				AllowedGeos:        entity.AllowedListType{},
				AllowedDevices:     entity.AllowedListType{},
				IsCampaignOveraged: false,
				IsCampaignActive:   false,
				CampaignDisabledRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.LinkRedirectType,
					RedirectURL:       "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignDisabled/CampaignDisableRedirectRulesLinkRedirectType",
					RedirectSlug:      "",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignDisabled/CampaignDisableRedirectRulesLinkRedirectType",
			expectedError:     nil,
		},
		{
			name: "CampaignDisableRedirectRulesLinkRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:                      true,
				AllowedProtocols:              entity.AllowedListType{},
				AllowedGeos:                   entity.AllowedListType{},
				AllowedDevices:                entity.AllowedListType{},
				IsCampaignOveraged:            false,
				IsCampaignActive:              false,
				CampaignDisabledRedirectRules: nil,
			},
			expectedTargetURL: "",
			expectedError:     interactor.ErrInvalidRedirectRules,
		},
		{
			name: "CampaignDisabledRedirectRulesSlugRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				AllowedProtocols:   entity.AllowedListType{},
				AllowedGeos:        entity.AllowedListType{},
				AllowedDevices:     entity.AllowedListType{},
				IsCampaignOveraged: false,
				IsCampaignActive:   false,
				CampaignDisabledRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.SlugRedirectType,
					RedirectURL:       "",
					RedirectSlug:      "testSlug123",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignDisabled/CampaignDisabledRedirectRulesSlugRedirectType",
			expectedError:     nil,
		},
		{
			name: "CampaignDisabledRedirectRulesSmartSlugRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				AllowedProtocols:   entity.AllowedListType{},
				AllowedGeos:        entity.AllowedListType{},
				AllowedDevices:     entity.AllowedListType{},
				IsCampaignOveraged: false,
				IsCampaignActive:   false,
				CampaignDisabledRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.SmartSlugRedirectType,
					RedirectURL:       "",
					RedirectSlug:      "",
					RedirectSmartSlug: []string{"testSlug000", "testSlug111", "testSlug222"},
				},
			},
			expectedTargetURL: "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignDisabled/CampaignDisabledRedirectRulesSmartSlugRedirectType",
			expectedError:     nil,
		},
		{
			name: "CampaignDisabledRedirectRulesNoRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				AllowedProtocols:   entity.AllowedListType{},
				AllowedGeos:        entity.AllowedListType{},
				AllowedDevices:     entity.AllowedListType{},
				IsCampaignOveraged: false,
				IsCampaignActive:   false,
				CampaignDisabledRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.NoRedirectType,
					RedirectURL:       "",
					RedirectSlug:      "",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "",
			expectedError:     interactor.ErrBlockRedirect,
		},
		{
			name: "CampaignDisabledRedirectRulesInvalidRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				AllowedProtocols:   entity.AllowedListType{},
				AllowedGeos:        entity.AllowedListType{},
				AllowedDevices:     entity.AllowedListType{},
				IsCampaignOveraged: false,
				IsCampaignActive:   false,
				CampaignDisabledRedirectRules: &valueobject.RedirectRules{
					RedirectType:      "invalid_redirect_type",
					RedirectURL:       "",
					RedirectSlug:      "",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "",
			expectedError:     interactor.ErrInvalidRedirectType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, srv, trkRepo, ipParser, uaParser, clkRepo := setupTest(t)
			defer ctrl.Finish()

			td := newTestData()
			tc.trkLink.Slug = td.slug

			trkRepo.EXPECT().FindTrackingLink(context.Background(), td.slug).Return(tc.trkLink)
			ipParser.EXPECT().Parse(td.requestData.IP).Return(td.countryCode, nil)
			uaParser.EXPECT().Parse(td.requestData.UserAgent).Return(td.userAgent, nil)

			if tc.expectedError == nil {
				clkRepo.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, click *entity.Click) error {
						if click.ID != td.requestData.RequestID {
							t.Errorf("unexpected click ID. expected %s but got %s", td.requestData.RequestID, click.ID)
						}
						if click.TargetURL != tc.expectedTargetURL {
							t.Errorf("unexpected target URL. expected %s but got %s", tc.expectedTargetURL, click.TargetURL)
						}

						return nil
					})
			}

			if tc.trkLink.CampaignDisabledRedirectRules != nil &&
				(tc.trkLink.CampaignDisabledRedirectRules.RedirectType == valueobject.SlugRedirectType ||
					tc.trkLink.CampaignDisabledRedirectRules.RedirectType == valueobject.SmartSlugRedirectType) {
				trkRepo.
					EXPECT().
					FindTrackingLink(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_, arg interface{}) *entity.TrackingLink {
						slug, ok := arg.(string)
						if !ok {
							t.Error("invalid argument type. expected string")
						}

						if tc.trkLink.CampaignDisabledRedirectRules.RedirectType == valueobject.SlugRedirectType {
							if slug != tc.trkLink.CampaignDisabledRedirectRules.RedirectSlug {
								t.Errorf("invalid argument received. expected %s but got %s", tc.trkLink.CampaignDisabledRedirectRules.RedirectSlug, slug)
							}
						} else if tc.trkLink.CampaignDisabledRedirectRules.RedirectType == valueobject.SmartSlugRedirectType {
							inArray := false

							for _, sl := range tc.trkLink.CampaignDisabledRedirectRules.RedirectSmartSlug {
								if sl == slug {
									inArray = true
								}
							}

							if !inArray {
								t.Errorf("invalid argument received. expected one of %v but got %s", tc.trkLink.CampaignDisabledRedirectRules.RedirectSmartSlug, slug)
							}
						}

						return &entity.TrackingLink{
							IsActive:           true,
							Slug:               slug,
							AllowedProtocols:   entity.AllowedListType{},
							AllowedGeos:        entity.AllowedListType{},
							AllowedDevices:     entity.AllowedListType{},
							IsCampaignOveraged: false,
							IsCampaignActive:   true,
							TargetURLTemplate:  "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignDisabled/" + tc.name,
						}
					})
				ipParser.EXPECT().Parse(td.requestData.IP).Return(td.countryCode, nil)
				uaParser.EXPECT().Parse(td.requestData.UserAgent).Return(td.userAgent, nil)
			}

			rResult, err := srv.Redirect(context.Background(), td.slug, td.requestData)
			if tc.expectedError != nil {
				if !errors.Is(err, tc.expectedError) {
					t.Errorf("unexpected result, %s expected", tc.expectedError.Error())
				}
				if rResult != nil {
					t.Error("unexpected target url. expected empty value")
				}
			} else if tc.expectedTargetURL != rResult.TargetURL {
				t.Errorf("unexpected target url. expected %s but got %s", tc.expectedTargetURL, rResult.TargetURL)
			} else {
				<-rResult.OutputCh
			}
		})
	}
}

func TestRedirectInteractor_Redirect_RenderTokens(t *testing.T) {
	trackingLink := entity.TrackingLink{
		IsActive:           true,
		AllowedProtocols:   entity.AllowedListType{},
		AllowedGeos:        entity.AllowedListType{},
		AllowedDevices:     entity.AllowedListType{},
		IsCampaignOveraged: false,
		IsCampaignActive:   true,
		TargetURLTemplate:  "http://target.url/path",
		CampaignID:         "1111",
		AffiliateID:        "1121",
		AdvertiserID:       "1131",
		SourceID:           "1141",
	}
	testCases := []struct {
		name              string
		trkLink           entity.TrackingLink
		tokens            []string
		expectedTargetURL string
	}{
		{
			name:              "RenderTokens_NoTokens",
			trkLink:           trackingLink,
			tokens:            []string{},
			expectedTargetURL: trackingLink.TargetURLTemplate,
		},
		{
			name:    "RenderTokens_IPAddressToken",
			trkLink: trackingLink,
			tokens:  []string{"{ip}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s",
				trackingLink.TargetURLTemplate, ipAddress),
		},
		{
			name:    "RenderTokens_ClickIDToken",
			trkLink: trackingLink,
			tokens:  []string{"{click_id}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s",
				trackingLink.TargetURLTemplate, requestID),
		},
		{
			name:    "RenderTokens_UserAgentToken",
			trkLink: trackingLink,
			tokens:  []string{"{user_agent}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s",
				trackingLink.TargetURLTemplate, userAgent),
		},
		{
			name:    "RenderTokens_CampaignIDToken",
			trkLink: trackingLink,
			tokens:  []string{"{campaign_id}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s",
				trackingLink.TargetURLTemplate, trackingLink.CampaignID),
		},
		{
			name:    "RenderTokens_AffiliateIDToken",
			trkLink: trackingLink,
			tokens:  []string{"{aff_id}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s",
				trackingLink.TargetURLTemplate, trackingLink.AffiliateID),
		},
		{
			name:    "RenderTokens_SourceIDToken",
			trkLink: trackingLink,
			tokens:  []string{"{source_id}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s",
				trackingLink.TargetURLTemplate, trackingLink.SourceID),
		},
		{
			name:    "RenderTokens_AdvertiserIDToken",
			trkLink: trackingLink,
			tokens:  []string{"{advertiser_id}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s",
				trackingLink.TargetURLTemplate, trackingLink.AdvertiserID),
		},
		{
			name:    "RenderTokens_DateToken",
			trkLink: trackingLink,
			tokens:  []string{"{date}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s",
				trackingLink.TargetURLTemplate, time.Now().Format("2006-01-02")),
		},
		{
			name:    "RenderTokens_DateTimeToken",
			trkLink: trackingLink,
			tokens:  []string{"{date_time}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s",
				trackingLink.TargetURLTemplate, time.Now().Format("2006-01-02T15:04:05")),
		},
		{
			name:    "RenderTokens_TimestampToken",
			trkLink: trackingLink,
			tokens:  []string{"{timestamp}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s",
				trackingLink.TargetURLTemplate, strconv.FormatInt(time.Now().Unix(), 10)),
		},
		{
			name:    "RenderTokens_P1-P4Tokens",
			trkLink: trackingLink,
			tokens:  []string{"{p1}", "{p2}", "{p3}", "{p4}"},
			expectedTargetURL: fmt.Sprintf(
				"%s?key0=%s&key1=%s&key2=%s&key3=%s",
				trackingLink.TargetURLTemplate, p1, p2, p3, p4),
		},
		{
			name:              "RenderTokens_CountryCodeToken",
			trkLink:           trackingLink,
			tokens:            []string{"{country_code}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s", trackingLink.TargetURLTemplate, countryCode),
		},
		{
			name:              "RenderTokens_RefererToken",
			trkLink:           trackingLink,
			tokens:            []string{"{referer}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s", trackingLink.TargetURLTemplate, referrer),
		},
		/*		{
					name:              "RenderTokens_RandomStrToken",
					trkLink:           expectedTrkLink,
					tokens:            []string{"{random_str}"},
					expectedTargetURL: fmt.Sprintf("%s?key0=%s", expectedTrkLink.TargetURLTemplate, ""), //TODO:
				},
				{
					name:              "RenderTokens_RandomIntToken",
					trkLink:           expectedTrkLink,
					tokens:            []string{"{random_int}"},
					expectedTargetURL: fmt.Sprintf("%s?key0=%s", expectedTrkLink.TargetURLTemplate, "1"), //TODO:
				},*/
		{
			name:              "RenderTokens_DeviceToken",
			trkLink:           trackingLink,
			tokens:            []string{"{device}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s", trackingLink.TargetURLTemplate, device),
		},
		{
			name:              "RenderTokens_PlatformToken",
			trkLink:           trackingLink,
			tokens:            []string{"{platform}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s", trackingLink.TargetURLTemplate, platform),
		},
		{
			name:              "RenderTokens_UnknownToken",
			trkLink:           trackingLink,
			tokens:            []string{"{unknown}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s", trackingLink.TargetURLTemplate, ""),
		},
		//TODO: test other tokens
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, srv, trkRepo, ipParser, uaParser, clkRepo := setupTest(t)
			defer ctrl.Finish()

			td := newTestData()
			tc.trkLink.Slug = td.slug

			if len(tc.tokens) > 0 {
				if !strings.HasSuffix(tc.trkLink.TargetURLTemplate, "?") {
					tc.trkLink.TargetURLTemplate += "?"
				}

				for i, token := range tc.tokens {
					tc.trkLink.TargetURLTemplate += "key" + strconv.Itoa(i) + "=" + token
					if i != (len(tc.tokens) - 1) {
						tc.trkLink.TargetURLTemplate += "&"
					}
				}
			}

			trkRepo.EXPECT().FindTrackingLink(gomock.Any(), td.slug).Return(&tc.trkLink)
			ipParser.EXPECT().Parse(td.requestData.IP).Return(td.countryCode, nil)
			uaParser.EXPECT().Parse(td.requestData.UserAgent).Return(td.userAgent, nil)

			clkRepo.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				// Return(nil)
				DoAndReturn(func(_ context.Context, click *entity.Click) error {
					if click.ID != td.requestData.RequestID {
						t.Errorf("unexpected click ID. expected %s but got %s", td.requestData.RequestID, click.ID)
					}
					if click.TargetURL != tc.expectedTargetURL {
						t.Errorf("unexpected target URL. expected %s but got %s", tc.expectedTargetURL, click.TargetURL)
					}

					return nil
				})

			rResult, err := srv.Redirect(context.Background(), td.slug, td.requestData)

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if rResult.TargetURL != tc.expectedTargetURL {
				t.Errorf("unexpected target URL. expected %s but got %s", tc.expectedTargetURL, rResult.TargetURL)
			}

			<-rResult.OutputCh
		})
	}
}

func TestRedirectInteractor_Redirect_LogErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv, trkRepo, ipAddressParser, userAgentParser := makeRedirectInteractor(ctrl)
	expectedSlug := "testSlug123"

	incomeURL, err := url.Parse("http://localhost/" + expectedSlug)
	if err != nil {
		t.Fatal(err)
	}

	expectedDto := &dto.RedirectRequestData{
		Params:    map[string][]string{"p1": []string{"val1"}, "p2": []string{"val2"}, "p4": []string{"val4"}},
		Headers:   make(map[string][]string),
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		IP:        net.ParseIP("178.43.146.107"),
		Protocol:  "http",
		Referer:   "https://httpbin.org",
		RequestID: "someUniqueRequestID",
		URL:       incomeURL,
	}
	expectedTrkLink := entity.TrackingLink{
		IsActive:           true,
		Slug:               expectedSlug,
		AllowedProtocols:   entity.AllowedListType{},
		AllowedGeos:        entity.AllowedListType{},
		AllowedDevices:     entity.AllowedListType{},
		IsCampaignOveraged: false,
		IsCampaignActive:   true,
		TargetURLTemplate:  "http://target.url/path",
		CampaignID:         "1234",
	}
	expectedCountry := ""
	expectedIPAddressParseError := errors.New("expected ip address parse error")
	expectedUserAgentParseError := errors.New("expected user agent parse error")
	expectedUserAgent := &valueobject.UserAgent{
		Bot:      false,
		Device:   "Mobile",
		Platform: "Android",
		Browser:  "Chrome",
	}

	trkRepo.EXPECT().FindTrackingLink(context.Background(), expectedSlug).Return(&expectedTrkLink)
	ipAddressParser.EXPECT().Parse(expectedDto.IP).Return(expectedCountry, expectedIPAddressParseError)
	userAgentParser.EXPECT().Parse(expectedDto.UserAgent).Return(expectedUserAgent, expectedUserAgentParseError)

	rResult, err := srv.Redirect(context.Background(), expectedSlug, expectedDto)
	if err != nil {
		t.Errorf("unexpected error. got %s\n", err)
	}
	if rResult.TargetURL != expectedTrkLink.TargetURLTemplate {
		t.Errorf("unexpected target url received. expected %s but got %s\n",
			expectedTrkLink.TargetURLTemplate, rResult.TargetURL)
	}
}

func TestRedirectInteractor_Redirect_OSValidation(t *testing.T) {
	tests := []struct {
		name           string
		trackingLink   *entity.TrackingLink
		userAgent      *valueobject.UserAgent
		expectedError  error
		expectedRules  *valueobject.RedirectRules
		shouldRedirect bool
	}{
		{
			name: "allowed OS platform",
			trackingLink: &entity.TrackingLink{
				IsCampaignOveraged: false,
				IsCampaignActive:   true,
				AllowedGeos:        make(entity.AllowedListType),
				AllowedDevices:     make(entity.AllowedListType),
				AllowedProtocols:   make(entity.AllowedListType),
				AllowDeeplink:      true,
				IsActive:           true,
				AllowedOS: entity.AllowedListType{
					"ios":     true,
					"android": true,
				},
			},
			userAgent: &valueobject.UserAgent{
				Platform: "ios",
			},
			shouldRedirect: true,
		},
		{
			name: "disallowed OS platform - redirect to specified URL",
			trackingLink: &entity.TrackingLink{
				IsCampaignOveraged: false,
				IsCampaignActive:   true,
				AllowedGeos:        make(entity.AllowedListType),
				AllowedDevices:     make(entity.AllowedListType),
				AllowedProtocols:   make(entity.AllowedListType),
				AllowDeeplink:      true,
				IsActive:           true,
				AllowedOS: entity.AllowedListType{
					"ios":     true,
					"android": true,
				},
				CampaignOSRedirectRules: &valueobject.RedirectRules{
					RedirectType: valueobject.LinkRedirectType,
					RedirectURL:  "https://os-restricted.example.com",
				},
			},
			userAgent: &valueobject.UserAgent{
				Platform: "windows",
			},
			expectedError:  nil,
			expectedRules:  &valueobject.RedirectRules{RedirectType: "url", RedirectURL: "https://os-restricted.example.com"},
			shouldRedirect: true,
		},
		{
			name: "disallowed OS platform",
			trackingLink: &entity.TrackingLink{
				IsCampaignOveraged: false,
				IsCampaignActive:   true,
				AllowedGeos:        make(entity.AllowedListType),
				AllowedDevices:     make(entity.AllowedListType),
				AllowedProtocols:   make(entity.AllowedListType),
				AllowDeeplink:      true,
				IsActive:           true,
				AllowedOS: entity.AllowedListType{
					"ios":     true,
					"android": true,
				},
				CampaignOSRedirectRules: &valueobject.RedirectRules{
					RedirectType: valueobject.NoRedirectType,
				},
			},
			userAgent: &valueobject.UserAgent{
				Platform: "windows",
			},
			expectedError:  interactor.ErrUnsupportedOS,
			expectedRules:  &valueobject.RedirectRules{RedirectType: valueobject.NoRedirectType},
			shouldRedirect: false,
		},
		{
			name: "disallowed OS platform - no redirect rules",
			trackingLink: &entity.TrackingLink{
				IsCampaignOveraged: false,
				IsCampaignActive:   true,
				AllowedGeos:        make(entity.AllowedListType),
				AllowedDevices:     make(entity.AllowedListType),
				AllowedProtocols:   make(entity.AllowedListType),
				AllowDeeplink:      true,
				IsActive:           true,
				AllowedOS: entity.AllowedListType{
					"ios":     true,
					"android": true,
				},
				CampaignOSRedirectRules: nil,
			},
			userAgent: &valueobject.UserAgent{
				Platform: "windows",
			},
			expectedError:  interactor.ErrUnsupportedOS,
			expectedRules:  &valueobject.RedirectRules{RedirectType: valueobject.NoRedirectType},
			shouldRedirect: false,
		},
		{
			name: "empty allowed OS list",
			trackingLink: &entity.TrackingLink{
				IsCampaignOveraged: false,
				IsCampaignActive:   true,
				AllowedGeos:        make(entity.AllowedListType),
				AllowedDevices:     make(entity.AllowedListType),
				AllowedProtocols:   make(entity.AllowedListType),
				AllowDeeplink:      true,
				IsActive:           true,
				AllowedOS:          entity.AllowedListType{},
			},
			userAgent: &valueobject.UserAgent{
				Platform: "windows",
			},
			shouldRedirect: true,
		},
		{
			name: "nil allowed OS list",
			trackingLink: &entity.TrackingLink{
				IsCampaignOveraged: false,
				IsCampaignActive:   true,
				AllowedGeos:        make(entity.AllowedListType),
				AllowedDevices:     make(entity.AllowedListType),
				AllowedProtocols:   make(entity.AllowedListType),
				AllowDeeplink:      true,
				IsActive:           true,
			},
			userAgent: &valueobject.UserAgent{
				Platform: "windows",
			},
			shouldRedirect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, srv, trkRepo, ipParser, uaParser, clkRepo := setupTest(t)
			defer ctrl.Finish()

			td := newTestData()

			// Configure mocks
			trkRepo.EXPECT().
				FindTrackingLink(gomock.Any(), gomock.Any()).
				Return(tt.trackingLink)

			uaParser.EXPECT().
				Parse(gomock.Any()).
				Return(tt.userAgent, nil)

			ipParser.EXPECT().
				Parse(gomock.Any()).
				Return(td.countryCode, nil)

			clkRepo.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()

			// Execute test
			result, err := srv.Redirect(context.Background(), td.slug, td.requestData)

			// Verify results
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.shouldRedirect && result == nil {
				t.Error("expected redirect result, got nil")
			}
		})
	}
}

// TestRedirectInteractor_Redirect_TokenReplacement tests the token replacement functionality
func TestRedirectInteractor_Redirect_TokenReplacement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		token             string
		expectedValue     string
		targetURLTemplate string
	}{
		{
			name:              "IP token",
			token:             "{ip}",
			expectedValue:     "192.168.1.1",
			targetURLTemplate: "https://example.com/track",
		},
		{
			name:              "user agent token",
			token:             "{user_agent}",
			expectedValue:     "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
			targetURLTemplate: "https://example.com/track",
		},
		{
			name:              "country code token",
			token:             "{country_code}",
			expectedValue:     "US",
			targetURLTemplate: "https://example.com/track",
		},
		{
			name:              "device token",
			token:             "{device}",
			expectedValue:     "Mobile",
			targetURLTemplate: "https://example.com/track",
		},
		{
			name:              "platform token",
			token:             "{platform}",
			expectedValue:     "Android",
			targetURLTemplate: "https://example.com/track",
		},
		{
			name:              "multiple tokens",
			token:             "{device}/{platform}/{country_code}",
			expectedValue:     "Mobile/Android/US",
			targetURLTemplate: "https://example.com/track",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, srv, trkRepo, ipParser, uaParser, clkRepo := setupTest(t)
			defer ctrl.Finish()

			td := newTestData()

			trkLink := &entity.TrackingLink{
				IsActive:           true,
				IsCampaignActive:   true,
				IsCampaignOveraged: false,
				TargetURLTemplate:  fmt.Sprintf("%s?value=%s", tt.targetURLTemplate, tt.token),
				AllowedProtocols:   make(entity.AllowedListType),
				AllowedGeos:        make(entity.AllowedListType),
				AllowedDevices:     make(entity.AllowedListType),
			}

			// Configure mocks
			trkRepo.EXPECT().
				FindTrackingLink(gomock.Any(), td.slug).
				Return(trkLink)

			uaParser.EXPECT().
				Parse(gomock.Any()).
				Return(td.userAgent, nil)

			ipParser.EXPECT().
				Parse(gomock.Any()).
				Return(td.countryCode, nil)

			clkRepo.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()

			// Execute test
			result, err := srv.Redirect(context.Background(), td.slug, td.requestData)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			expectedURL := fmt.Sprintf("%s?value=%s", tt.targetURLTemplate, tt.expectedValue)
			if result.TargetURL != expectedURL {
				t.Errorf("expected URL %s, got %s", expectedURL, result.TargetURL)
			}
		})
	}
}
