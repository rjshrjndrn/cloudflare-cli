package cloudflare

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		email     string
		wantError bool
	}{
		{
			name:      "API Token only",
			token:     "test-token",
			email:     "",
			wantError: false,
		},
		{
			name:      "API Key with email",
			token:     "test-key",
			email:     "test@example.com",
			wantError: false,
		},
		{
			name:      "Empty token",
			token:     "",
			email:     "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.token, tt.email)
			if (err != nil) != tt.wantError {
				t.Errorf("NewClient() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && client == nil {
				t.Error("Expected client to be non-nil")
			}
		})
	}
}
