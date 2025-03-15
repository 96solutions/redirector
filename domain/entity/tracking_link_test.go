package entity_test

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/lroman242/redirector/domain/entity"
)

func TestAllowedListType_Value(t *testing.T) {
	tests := []struct {
		name    string
		list    entity.AllowedListType
		want    driver.Value
		wantErr bool
	}{
		{
			name: "empty map",
			list: entity.AllowedListType{},
			want: "{}",
		},
		{
			name: "single value",
			list: entity.AllowedListType{"http": true},
			want: `{"http":true}`,
		},
		{
			name: "multiple values",
			list: entity.AllowedListType{
				"http":  true,
				"https": false,
			},
			want: `{"http":true,"https":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.list.Value()

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("AllowedListType.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Convert to expected format for comparison
			if gotBytes, ok := got.([]byte); ok {
				gotStr := string(gotBytes)
				wantStr := tt.want.(string)

				// Convert both to maps to handle JSON key order differences
				var gotMap, wantMap map[string]bool
				if err := json.Unmarshal([]byte(gotStr), &gotMap); err != nil {
					t.Errorf("Failed to unmarshal got value: %v", err)
					return
				}
				if err := json.Unmarshal([]byte(wantStr), &wantMap); err != nil {
					t.Errorf("Failed to unmarshal want value: %v", err)
					return
				}

				if !reflect.DeepEqual(gotMap, wantMap) {
					t.Errorf("AllowedListType.Value() = %v, want %v", gotStr, wantStr)
				}
			} else {
				t.Errorf("AllowedListType.Value() returned unexpected type: %T", got)
			}
		})
	}
}

func TestAllowedListType_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    entity.AllowedListType
		wantErr bool
	}{
		{
			name:    "nil value",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "non-byte slice value",
			input:   "string value",
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{"invalid json`),
			wantErr: true,
		},
		{
			name:  "empty JSON object",
			input: []byte(`{}`),
			want:  entity.AllowedListType{},
		},
		{
			name:  "single value",
			input: []byte(`{"http":true}`),
			want:  entity.AllowedListType{"http": true},
		},
		{
			name:  "multiple values",
			input: []byte(`{"http":true,"https":false}`),
			want:  entity.AllowedListType{"http": true, "https": false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got entity.AllowedListType
			err := got.Scan(tt.input)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("AllowedListType.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If expecting error, don't compare results
			if tt.wantErr {
				return
			}

			// Check result
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllowedListType.Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllowedListType_Scan_TypeAssertionError(t *testing.T) {
	var list entity.AllowedListType
	err := list.Scan(123) // Pass an integer instead of a byte slice

	if err == nil {
		t.Errorf("Expected an error for type assertion failure, got nil")
	}

	expectedErr := "type assertion to []byte failed"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message %q, got %q", expectedErr, err.Error())
	}
}

func TestAllowedListType_ScanUnmarshalError(t *testing.T) {
	var list entity.AllowedListType
	invalidJSON := []byte(`{"invalid":json}`)

	err := list.Scan(invalidJSON)

	if err == nil {
		t.Errorf("Expected a JSON unmarshaling error, got nil")
	}

	// Check that it's a JSON syntax error
	var jsonErr *json.SyntaxError
	if !errors.As(err, &jsonErr) {
		t.Errorf("Expected a JSON syntax error, got %T: %v", err, err)
	}
}
