package dto_test

import (
	"github.com/lroman242/redirector/domain/dto"
	"strings"
	"testing"
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
