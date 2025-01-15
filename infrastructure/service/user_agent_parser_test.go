package service_test

import (
	"errors"
	"testing"

	"github.com/lroman242/redirector/domain/valueobject"
	"github.com/lroman242/redirector/infrastructure/service"
)

func TestUserAgentParser_Parse(t *testing.T) {
	testCases := []struct {
		name                 string
		userAgentStr         string
		expectedUserAgentObj *valueobject.UserAgent
		expectedError        error
	}{
		{
			name:         "Mac|MacOS|Safari",
			userAgentStr: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
			expectedUserAgentObj: &valueobject.UserAgent{
				SrcString: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
				Bot:       false,
				Device:    "Mac",
				Platform:  "Mac OS X",
				Browser:   "Safari",
			},
			expectedError: nil,
		},
		{
			name:         "Samsung|Android|Opera",
			userAgentStr: "Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956",
			expectedUserAgentObj: &valueobject.UserAgent{
				SrcString: "Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956",
				Bot:       false,
				Device:    "Samsung GT-I9300",
				Platform:  "Android",
				Browser:   "Opera Mobile",
			},
			expectedError: nil,
		},
		{
			name:         "empty user agent",
			userAgentStr: "",
			expectedUserAgentObj: &valueobject.UserAgent{
				SrcString: "",
				Bot:       false,
				Device:    "",
				Platform:  "",
				Browser:   "",
			},
			expectedError: service.EmptyUserAgentError,
		},
	}

	parser := service.NewUserAgentParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ua, err := parser.Parse(tc.userAgentStr)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("unexpected error: %s\n", err)
			}

			if ua.Browser != tc.expectedUserAgentObj.Browser {
				t.Errorf("unexpected browser value parsed. expected %s but got %s\n",
					tc.expectedUserAgentObj.Browser, ua.Browser)
			}
			if ua.Platform != tc.expectedUserAgentObj.Platform {
				t.Errorf("unexpected platform value parsed. expected %s but got %s\n",
					tc.expectedUserAgentObj.Platform, ua.Platform)
			}
			if ua.Device != tc.expectedUserAgentObj.Device {
				t.Errorf("unexpected device value parsed. expected %s but got %s\n",
					tc.expectedUserAgentObj.Device, ua.Device)
			}
		})
	}
}
