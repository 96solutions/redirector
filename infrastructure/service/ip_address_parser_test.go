package service_test

import (
	"errors"
	"net"
	"testing"

	"github.com/lroman242/redirector/infrastructure/service"
	"github.com/oschwald/geoip2-golang"
)

func TestGeoIP2_Parse(t *testing.T) {
	reader, err := geoip2.Open("./../../GeoLite2-Country.mmdb")
	if err != nil {
		t.Errorf("cannot initialize GeoIP parser")
	}

	geoIpParser := service.NewGeoIP2(reader)

	testCases := []struct {
		name            string
		ip              string
		expectedCountry string
		expectedError   error
	}{
		{
			name:            "PL ip address",
			ip:              "178.43.70.56",
			expectedCountry: "PL",
			expectedError:   nil,
		},
		{
			name:            "SG ip address",
			ip:              "206.189.156.75",
			expectedCountry: "SG",
			expectedError:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			countryCode, internalErr := geoIpParser.Parse(net.ParseIP(tc.ip))
			if !errors.Is(internalErr, tc.expectedError) {
				t.Errorf("unexpected error. expected %s but got %s\n", tc.expectedError, err)
			}
			if countryCode != tc.expectedCountry {
				t.Errorf("unexpected country code. expected %s but got %s\n", tc.expectedCountry, countryCode)
			}
		})
	}
}
