package http

import (
	"net"
	"net/http"
	"testing"
)

func TestGetIPAddress(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		expectedIP net.IP
		expectErr  bool
	}{
		{
			name: "CF-Connecting-IP header",
			headers: map[string]string{
				"CF-Connecting-IP": "203.0.113.1",
			},
			expectedIP: net.ParseIP("203.0.113.1"),
			expectErr:  false,
		},
		{
			name: "X-Real-IP header",
			headers: map[string]string{
				"X-Real-IP": "203.0.113.2",
			},
			expectedIP: net.ParseIP("203.0.113.2"),
			expectErr:  false,
		},
		{
			name: "X-Forwarded-For header",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.3, 203.0.113.4",
			},
			expectedIP: net.ParseIP("203.0.113.3"),
			expectErr:  false,
		},
		{
			name:       "RemoteAddr fallback",
			remoteAddr: "203.0.113.5:12345",
			expectedIP: net.ParseIP("203.0.113.5"),
			expectErr:  false,
		},
		{
			name:       "Invalid RemoteAddr",
			remoteAddr: "invalid-addr",
			expectedIP: nil,
			expectErr:  true,
		},
		{
			name:       "No headers and invalid RemoteAddr",
			remoteAddr: "",
			expectedIP: nil,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header:     make(http.Header),
				RemoteAddr: tt.remoteAddr,
			}
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			ip, err := getIPAddress(req)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error", tt.expectErr)
				}
			} else {
				if err != nil {
					t.Error("unexpected error", err)
				}
				if tt.expectedIP.String() != ip.String() {
					t.Errorf("unexpected result. expected %s got %s", tt.expectedIP.String(), ip.String())
				}
			}
		})
	}
}
