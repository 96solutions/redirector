package interactor

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/lroman242/redirector/domain/dto"
	"github.com/lroman242/redirector/domain/entity"
	"github.com/lroman242/redirector/domain/valueobject"
	"github.com/lroman242/redirector/mocks"
	"net"
	"strconv"
	"strings"
	"testing"
)

func TestRedirectInteractor_Redirect_TrackingLinkNotFoundError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedDto := &dto.RedirectRequestData{
		Params:    make(map[string][]string),
		Headers:   make(map[string]string),
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		IP:        net.ParseIP("178.43.146.107"),
		Protocol:  "http",
	}
	expectedSlug := "testSlug123"

	repo := mocks.NewMockTrackingLinksRepositoryInterface(ctrl)
	repo.EXPECT().FindTrackingLink(expectedSlug).Return(nil)

	ipAddressParser := mocks.NewMockIpAddressParserInterface(ctrl)
	userAgentParser := mocks.NewMockUserAgentParser(ctrl)

	srv := NewRedirectInteractor(repo, ipAddressParser, userAgentParser)
	targetURL, err := srv.Redirect(context.Background(), expectedSlug, expectedDto)

	if !errors.Is(err, TrackingLinkNotFoundError) {
		t.Error("unexpected result, TrackingLinkNotFoundError expected")
	}
	if targetURL != "" {
		t.Error("unexpected target url. expected empty value")
	}
}

func TestRedirectInteractor_Redirect_DisabledTrackingLinkError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockTrackingLinksRepositoryInterface(ctrl)
	ipAddressParser := mocks.NewMockIpAddressParserInterface(ctrl)
	userAgentParser := mocks.NewMockUserAgentParser(ctrl)

	srv := NewRedirectInteractor(repo, ipAddressParser, userAgentParser)

	expectedDto := &dto.RedirectRequestData{
		Params:    make(map[string][]string),
		Headers:   make(map[string]string),
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		IP:        net.ParseIP("178.43.146.107"),
		Protocol:  "http",
	}
	expectedSlug := "testSlug123"

	trkLink := &entity.TrackingLink{
		IsActive:         false,
		Slug:             expectedSlug,
		AllowedProtocols: []string{"https"},
	}

	repo.EXPECT().FindTrackingLink(expectedSlug).Return(trkLink)

	targetURL, err := srv.Redirect(context.Background(), expectedSlug, expectedDto)
	if !errors.Is(err, TrackingLinkDisabledError) {
		t.Error("unexpected result, TrackingLinkDisabledError expected")
	}
	if targetURL != "" {
		t.Error("unexpected target url. expected empty value")
	}
}

func TestRedirectInteractor_Redirect_WrongProtocolError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockTrackingLinksRepositoryInterface(ctrl)
	ipAddressParser := mocks.NewMockIpAddressParserInterface(ctrl)
	userAgentParser := mocks.NewMockUserAgentParser(ctrl)

	srv := NewRedirectInteractor(repo, ipAddressParser, userAgentParser)

	expectedDto := &dto.RedirectRequestData{
		Params:    make(map[string][]string),
		Headers:   make(map[string]string),
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		IP:        net.ParseIP("178.43.146.107"),
		Protocol:  "http",
	}
	expectedSlug := "testSlug123"

	trkLink := &entity.TrackingLink{
		IsActive:         true,
		Slug:             expectedSlug,
		AllowedProtocols: []string{"https"},
	}

	repo.EXPECT().FindTrackingLink(expectedSlug).Return(trkLink)

	targetURL, err := srv.Redirect(context.Background(), expectedSlug, expectedDto)
	if !errors.Is(err, UnsupportedProtocolError) {
		t.Error("unexpected result, UnsupportedProtocolError expected")
	}
	if targetURL != "" {
		t.Error("unexpected target url. expected empty value")
	}
}

func TestRedirectInteractor_Redirect_WrongGeoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockTrackingLinksRepositoryInterface(ctrl)
	ipAddressParser := mocks.NewMockIpAddressParserInterface(ctrl)
	userAgentParser := mocks.NewMockUserAgentParser(ctrl)

	srv := NewRedirectInteractor(repo, ipAddressParser, userAgentParser)

	expectedDto := &dto.RedirectRequestData{
		Params:    make(map[string][]string),
		Headers:   make(map[string]string),
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		IP:        net.ParseIP("178.43.146.107"),
		Protocol:  "http",
	}
	expectedSlug := "testSlug123"

	trkLink := &entity.TrackingLink{
		IsActive:         true,
		Slug:             expectedSlug,
		AllowedProtocols: []string{},
		AllowedGeos:      []string{"US", "PT", "UA"},
	}

	repo.EXPECT().FindTrackingLink(expectedSlug).Return(trkLink)
	ipAddressParser.EXPECT().Parse(expectedDto.IP).Return("PL", nil)

	targetURL, err := srv.Redirect(context.Background(), expectedSlug, expectedDto)
	if !errors.Is(err, UnsupportedGeoError) {
		t.Error("unexpected result, UnsupportedGeoError expected")
	}
	if targetURL != "" {
		t.Error("unexpected target url. expected empty value")
	}
}

func TestRedirectInteractor_Redirect_WrongDeviceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockTrackingLinksRepositoryInterface(ctrl)
	ipAddressParser := mocks.NewMockIpAddressParserInterface(ctrl)
	userAgentParser := mocks.NewMockUserAgentParser(ctrl)

	srv := NewRedirectInteractor(repo, ipAddressParser, userAgentParser)

	expectedDto := &dto.RedirectRequestData{
		Params:    make(map[string][]string),
		Headers:   make(map[string]string),
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		IP:        net.ParseIP("178.43.146.107"),
		Protocol:  "http",
	}
	expectedSlug := "testSlug123"

	trkLink := &entity.TrackingLink{
		IsActive:         true,
		Slug:             expectedSlug,
		AllowedProtocols: []string{},
		AllowedGeos:      []string{"US", "PT", "UA", "PL"},
		AllowedDevices:   []string{"Desktop"},
	}

	repo.EXPECT().FindTrackingLink(expectedSlug).Return(trkLink)
	ipAddressParser.EXPECT().Parse(expectedDto.IP).Return("PL", nil)
	userAgentParser.EXPECT().Parse(expectedDto.UserAgent).Return(&valueobject.UserAgent{
		Bot:      false,
		Device:   "Mobile",
		Platform: "Android",
		Browser:  "Chrome",
		Version:  "109.0.5414.119",
	}, nil)

	targetURL, err := srv.Redirect(context.Background(), expectedSlug, expectedDto)
	if !errors.Is(err, UnsupportedDeviceError) {
		t.Error("unexpected result, UnsupportedDeviceError expected")
	}
	if targetURL != "" {
		t.Error("unexpected target url. expected empty value")
	}
}

func TestRedirectInteractor_Redirect_CampaignOveraged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockTrackingLinksRepositoryInterface(ctrl)
	ipAddressParser := mocks.NewMockIpAddressParserInterface(ctrl)
	userAgentParser := mocks.NewMockUserAgentParser(ctrl)

	srv := NewRedirectInteractor(repo, ipAddressParser, userAgentParser)

	expectedDto := &dto.RedirectRequestData{
		Params:    make(map[string][]string),
		Headers:   make(map[string]string),
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		IP:        net.ParseIP("178.43.146.107"),
		Protocol:  "http",
	}
	expectedSlug := "testSlug123"
	expectedSlug2 := "testSlug456"
	expectedSlug3 := []string{"testSlug000", "testSlug111", "testSlug222"}

	expectedTargetURL := "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignOveraged"

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
				Slug:               expectedSlug,
				AllowedProtocols:   []string{},
				AllowedGeos:        []string{},
				AllowedDevices:     []string{},
				IsCampaignOveraged: true,
				CampaignOverageRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.LinkRedirectType,
					RedirectURL:       "http://sometarget.url/test",
					RedirectSlug:      "",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "http://sometarget.url/test",
			expectedError:     nil,
		},
		{
			name: "OveragedRedirectRulesSlugRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				Slug:               expectedSlug,
				AllowedProtocols:   []string{},
				AllowedGeos:        []string{},
				AllowedDevices:     []string{},
				IsCampaignOveraged: true,
				CampaignOverageRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.SlugRedirectType,
					RedirectURL:       "",
					RedirectSlug:      expectedSlug2,
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: expectedTargetURL,
			expectedError:     nil,
		},
		{
			name: "OveragedRedirectRulesSmartSlugRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				Slug:               expectedSlug,
				AllowedProtocols:   []string{},
				AllowedGeos:        []string{},
				AllowedDevices:     []string{},
				IsCampaignOveraged: true,
				CampaignOverageRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.SmartSlugRedirectType,
					RedirectURL:       "",
					RedirectSlug:      "",
					RedirectSmartSlug: expectedSlug3,
				},
			},
			expectedTargetURL: expectedTargetURL,
			expectedError:     nil,
		},
		{
			name: "OveragedRedirectRulesNoRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				Slug:               expectedSlug,
				AllowedProtocols:   []string{},
				AllowedGeos:        []string{},
				AllowedDevices:     []string{},
				IsCampaignOveraged: true,
				CampaignOverageRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.NoRedirectType,
					RedirectURL:       "",
					RedirectSlug:      "",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "",
			expectedError:     BlockRedirectError,
		},
		{
			name: "OveragedRedirectRulesInvalidRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				Slug:               expectedSlug,
				AllowedProtocols:   []string{},
				AllowedGeos:        []string{},
				AllowedDevices:     []string{},
				IsCampaignOveraged: true,
				CampaignOverageRedirectRules: &valueobject.RedirectRules{
					RedirectType:      "invalid_redirect_type",
					RedirectURL:       "",
					RedirectSlug:      "",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "",
			expectedError:     InvalidRedirectTypeError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo.EXPECT().FindTrackingLink(expectedSlug).Return(tc.trkLink)
			ipAddressParser.EXPECT().Parse(expectedDto.IP).Return("PL", nil)
			userAgentParser.EXPECT().Parse(expectedDto.UserAgent).Return(&valueobject.UserAgent{
				Bot:      false,
				Device:   "Mobile",
				Platform: "Android",
				Browser:  "Chrome",
				Version:  "109.0.5414.119",
			}, nil)

			if tc.trkLink.CampaignOverageRedirectRules.RedirectType == valueobject.SlugRedirectType ||
				tc.trkLink.CampaignOverageRedirectRules.RedirectType == valueobject.SmartSlugRedirectType {
				repo.EXPECT().FindTrackingLink(gomock.Any()).DoAndReturn(func(arg interface{}) *entity.TrackingLink {
					if tc.trkLink.CampaignOverageRedirectRules.RedirectType == valueobject.SlugRedirectType {
						slug, ok := arg.(string)
						if !ok {
							t.Error("invalid argument type. expected string")
						}
						if slug != expectedSlug2 {
							t.Errorf("invalid argument received. expected %s but got %s", expectedSlug2, slug)
						}
					} else if tc.trkLink.CampaignOverageRedirectRules.RedirectType == valueobject.SmartSlugRedirectType {
						slug, ok := arg.(string)
						if !ok {
							t.Error("invalid argument type. expected string")
						}

						inArray := false

						for _, sl := range expectedSlug3 {
							if sl == slug {
								inArray = true
							}
						}

						if !inArray {
							t.Errorf("invalid argument received. expected one of %v but got %s", expectedSlug3, slug)
						}
					}

					return &entity.TrackingLink{
						IsActive:           true,
						Slug:               expectedSlug,
						AllowedProtocols:   []string{},
						AllowedGeos:        []string{},
						AllowedDevices:     []string{},
						IsCampaignOveraged: false,
						IsCampaignActive:   true,
						TargetURLTemplate:  expectedTargetURL,
					}
				})
				ipAddressParser.EXPECT().Parse(expectedDto.IP).Return("PL", nil)
				userAgentParser.EXPECT().Parse(expectedDto.UserAgent).Return(&valueobject.UserAgent{
					Bot:      false,
					Device:   "Mobile",
					Platform: "Android",
					Browser:  "Chrome",
					Version:  "109.0.5414.119",
				}, nil)
			}

			targetURL, err := srv.Redirect(context.Background(), expectedSlug, expectedDto)
			if tc.expectedError != nil {
				if !errors.Is(err, tc.expectedError) {
					t.Errorf("unexpected result, %T expected", tc.expectedError)
				}
				if targetURL != "" {
					t.Error("unexpected target url. expected empty value")
				}
			} else if tc.expectedTargetURL != targetURL {
				t.Errorf("unexpected target url. expected %s but got %s", tc.expectedTargetURL, targetURL)
			}

		})
	}

}

func TestRedirectInteractor_Redirect_CampaignDisabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockTrackingLinksRepositoryInterface(ctrl)
	ipAddressParser := mocks.NewMockIpAddressParserInterface(ctrl)
	userAgentParser := mocks.NewMockUserAgentParser(ctrl)

	srv := NewRedirectInteractor(repo, ipAddressParser, userAgentParser)

	expectedDto := &dto.RedirectRequestData{
		Params:    make(map[string][]string),
		Headers:   make(map[string]string),
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		IP:        net.ParseIP("178.43.146.107"),
		Protocol:  "http",
	}
	expectedSlug := "testSlug123"
	expectedSlug2 := "testSlug456"
	expectedSlug3 := []string{"testSlug000", "testSlug111", "testSlug222"}

	expectedTargetURL := "http://sometarget.url/TestRedirectInteractor_Redirect_CampaignOveraged"

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
				Slug:               expectedSlug,
				AllowedProtocols:   []string{},
				AllowedGeos:        []string{},
				AllowedDevices:     []string{},
				IsCampaignOveraged: false,
				IsCampaignActive:   false,
				CampaignDisabledRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.LinkRedirectType,
					RedirectURL:       "http://sometarget.url/test",
					RedirectSlug:      "",
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: "http://sometarget.url/test",
			expectedError:     nil,
		},
		{
			name: "CampaignDisabledRedirectRulesSlugRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				Slug:               expectedSlug,
				AllowedProtocols:   []string{},
				AllowedGeos:        []string{},
				AllowedDevices:     []string{},
				IsCampaignOveraged: false,
				IsCampaignActive:   false,
				CampaignDisabledRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.SlugRedirectType,
					RedirectURL:       "",
					RedirectSlug:      expectedSlug2,
					RedirectSmartSlug: nil,
				},
			},
			expectedTargetURL: expectedTargetURL,
			expectedError:     nil,
		},
		{
			name: "CampaignDisabledRedirectRulesSmartSlugRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				Slug:               expectedSlug,
				AllowedProtocols:   []string{},
				AllowedGeos:        []string{},
				AllowedDevices:     []string{},
				IsCampaignOveraged: false,
				IsCampaignActive:   false,
				CampaignDisabledRedirectRules: &valueobject.RedirectRules{
					RedirectType:      valueobject.SmartSlugRedirectType,
					RedirectURL:       "",
					RedirectSlug:      "",
					RedirectSmartSlug: expectedSlug3,
				},
			},
			expectedTargetURL: expectedTargetURL,
			expectedError:     nil,
		},
		{
			name: "CampaignDisabledRedirectRulesNoRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				Slug:               expectedSlug,
				AllowedProtocols:   []string{},
				AllowedGeos:        []string{},
				AllowedDevices:     []string{},
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
			expectedError:     BlockRedirectError,
		},
		{
			name: "CampaignDisabledRedirectRulesInvalidRedirectType",
			trkLink: &entity.TrackingLink{
				IsActive:           true,
				Slug:               expectedSlug,
				AllowedProtocols:   []string{},
				AllowedGeos:        []string{},
				AllowedDevices:     []string{},
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
			expectedError:     InvalidRedirectTypeError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo.EXPECT().FindTrackingLink(expectedSlug).Return(tc.trkLink)
			ipAddressParser.EXPECT().Parse(expectedDto.IP).Return("PL", nil)
			userAgentParser.EXPECT().Parse(expectedDto.UserAgent).Return(&valueobject.UserAgent{
				Bot:      false,
				Device:   "Mobile",
				Platform: "Android",
				Browser:  "Chrome",
				Version:  "109.0.5414.119",
			}, nil)

			if tc.trkLink.CampaignDisabledRedirectRules.RedirectType == valueobject.SlugRedirectType ||
				tc.trkLink.CampaignDisabledRedirectRules.RedirectType == valueobject.SmartSlugRedirectType {
				repo.EXPECT().FindTrackingLink(gomock.Any()).DoAndReturn(func(arg interface{}) *entity.TrackingLink {
					if tc.trkLink.CampaignDisabledRedirectRules.RedirectType == valueobject.SlugRedirectType {
						slug, ok := arg.(string)
						if !ok {
							t.Error("invalid argument type. expected string")
						}
						if slug != expectedSlug2 {
							t.Errorf("invalid argument received. expected %s but got %s", expectedSlug2, slug)
						}
					} else if tc.trkLink.CampaignDisabledRedirectRules.RedirectType == valueobject.SmartSlugRedirectType {
						slug, ok := arg.(string)
						if !ok {
							t.Error("invalid argument type. expected string")
						}

						inArray := false

						for _, sl := range expectedSlug3 {
							if sl == slug {
								inArray = true
							}
						}

						if !inArray {
							t.Errorf("invalid argument received. expected one of %v but got %s", expectedSlug3, slug)
						}
					}

					return &entity.TrackingLink{
						IsActive:           true,
						Slug:               expectedSlug,
						AllowedProtocols:   []string{},
						AllowedGeos:        []string{},
						AllowedDevices:     []string{},
						IsCampaignOveraged: false,
						IsCampaignActive:   true,
						TargetURLTemplate:  expectedTargetURL,
					}
				})
				ipAddressParser.EXPECT().Parse(expectedDto.IP).Return("PL", nil)
				userAgentParser.EXPECT().Parse(expectedDto.UserAgent).Return(&valueobject.UserAgent{
					Bot:      false,
					Device:   "Mobile",
					Platform: "Android",
					Browser:  "Chrome",
					Version:  "109.0.5414.119",
				}, nil)
			}

			targetURL, err := srv.Redirect(context.Background(), expectedSlug, expectedDto)
			if tc.expectedError != nil {
				if !errors.Is(err, tc.expectedError) {
					t.Errorf("unexpected result, %T expected", tc.expectedError)
				}
				if targetURL != "" {
					t.Error("unexpected target url. expected empty value")
				}
			} else if tc.expectedTargetURL != targetURL {
				t.Errorf("unexpected target url. expected %s but got %s", tc.expectedTargetURL, targetURL)
			}

		})
	}
}

func TestRedirectInteractor_Redirect_RenderTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockTrackingLinksRepositoryInterface(ctrl)
	ipAddressParser := mocks.NewMockIpAddressParserInterface(ctrl)
	userAgentParser := mocks.NewMockUserAgentParser(ctrl)

	expectedDto := &dto.RedirectRequestData{
		Params:    make(map[string][]string),
		Headers:   make(map[string]string),
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		IP:        net.ParseIP("178.43.146.107"),
		Protocol:  "http",
	}
	expectedSlug := "testSlug123"
	expectedTrkLink := entity.TrackingLink{
		IsActive:           true,
		Slug:               expectedSlug,
		AllowedProtocols:   []string{},
		AllowedGeos:        []string{},
		AllowedDevices:     []string{},
		IsCampaignOveraged: false,
		IsCampaignActive:   true,
		TargetURLTemplate:  "http://target.url/path",
	}
	expectedCountry := "PL"
	expectedUserAgent := &valueobject.UserAgent{
		Bot:      false,
		Device:   "Mobile",
		Platform: "Android",
		Browser:  "Chrome",
		Version:  "109.0.5414.119",
	}

	testCases := []struct {
		name              string
		trkLink           entity.TrackingLink
		tokens            []string
		expectedTargetURL string
	}{
		{
			name:              "RenderTokens_NoTokens",
			trkLink:           expectedTrkLink,
			tokens:            []string{},
			expectedTargetURL: expectedTrkLink.TargetURLTemplate,
		},
		{
			name:              "RenderTokens_IP",
			trkLink:           expectedTrkLink,
			tokens:            []string{"{ip}"},
			expectedTargetURL: fmt.Sprintf("%s?key0=%s", expectedTrkLink.TargetURLTemplate, expectedDto.IP),
		},
		//TODO: test other tokens
	}
	srv := NewRedirectInteractor(repo, ipAddressParser, userAgentParser)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

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

			repo.EXPECT().FindTrackingLink(expectedSlug).Return(&tc.trkLink)
			ipAddressParser.EXPECT().Parse(expectedDto.IP).Return(expectedCountry, nil)
			userAgentParser.EXPECT().Parse(expectedDto.UserAgent).Return(expectedUserAgent, nil)

			targetURL, err := srv.Redirect(context.Background(), expectedSlug, expectedDto)

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if targetURL != tc.expectedTargetURL {
				t.Errorf("unexpected target URL. expected %s but got %s", tc.expectedTargetURL, targetURL)
			}
		})
	}
}
