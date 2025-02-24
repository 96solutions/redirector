package dto_test

import (
	"strings"
	"testing"

	"net"

	"github.com/lroman242/redirector/domain/dto"
)

func TestRedirectRequestData_GetParam(t *testing.T) {
	testCases := []struct {
		name     string
		params   map[string][]string
		expected string
	}{
		{
			name:     "no param",
			params:   map[string][]string{},
			expected: "",
		},
		{
			name:     "single value per key",
			params:   map[string][]string{"key1": []string{"someValue1"}},
			expected: "someValue1",
		},
		{
			name:     "multiple values per key",
			params:   map[string][]string{"key1": []string{"someValue0", "someValue1", "someValue2"}},
			expected: "someValue0,someValue1,someValue2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requestData := dto.RedirectRequestData{
				Params: tc.params,
			}

			result := requestData.GetParam("key1")
			sResult := strings.Join(result, ",")
			if sResult != tc.expected {
				t.Errorf("unexpected result received. expected %s but got %s\n", tc.expected, sResult)
			}
		})
	}
}

func TestRedirectRequestData_Validate(t *testing.T) {
	tests := []struct {
		name        string
		data        *dto.RedirectRequestData
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid request data",
			data: &dto.RedirectRequestData{
				Slug: "test-slug",
				IP:   net.ParseIP("192.168.1.1"),
			},
			wantErr: false,
		},
		{
			name: "missing slug",
			data: &dto.RedirectRequestData{
				IP: net.ParseIP("192.168.1.1"),
			},
			wantErr:     true,
			expectedErr: "slug is required",
		},
		{
			name: "missing IP",
			data: &dto.RedirectRequestData{
				Slug: "test-slug",
			},
			wantErr:     true,
			expectedErr: "IP is required",
		},
		{
			name: "empty slug",
			data: &dto.RedirectRequestData{
				Slug: "",
				IP:   net.ParseIP("192.168.1.1"),
			},
			wantErr:     true,
			expectedErr: "slug is required",
		},
		{
			name: "invalid IP",
			data: &dto.RedirectRequestData{
				Slug: "test-slug",
				IP:   nil,
			},
			wantErr:     true,
			expectedErr: "IP is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()

			// Check if error was expected
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If error was expected, check error message
			if tt.wantErr && err.Error() != tt.expectedErr {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.expectedErr)
			}
		})
	}
}
