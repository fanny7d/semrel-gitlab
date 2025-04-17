package gitlabutil

import (
	"net/http"
	"testing"

	"github.com/xanzy/go-gitlab"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		apiURL      string
		skipSSL     bool
		expectError bool
	}{
		{
			name:        "Valid configuration",
			token:       "valid-token",
			apiURL:      "https://gitlab.example.com/api/v4",
			skipSSL:     false,
			expectError: false,
		},
		{
			name:        "Empty token",
			token:       "",
			apiURL:      "https://gitlab.example.com/api/v4",
			skipSSL:     false,
			expectError: true,
		},
		{
			name:        "Invalid API URL",
			token:       "valid-token",
			apiURL:      "not-a-url",
			skipSSL:     false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.token, tt.apiURL, tt.skipSSL)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if client == nil {
				t.Error("Expected client, got nil")
			}
		})
	}
}

func TestNewRequest(t *testing.T) {
	client, err := NewClient("test-token", "https://gitlab.example.com/api/v4", false)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name        string
		method      string
		path        string
		body        interface{}
		options     []gitlab.RequestOptionFunc
		expectError bool
	}{
		{
			name:        "Valid request",
			method:      http.MethodGet,
			path:        "/projects",
			body:        nil,
			options:     nil,
			expectError: false,
		},
		{
			name:        "Invalid method",
			method:      "INVALID",
			path:        "/projects",
			body:        nil,
			options:     nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := client.NewRequest(tt.method, tt.path, tt.body, tt.options)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if req == nil {
				t.Error("Expected request, got nil")
			}

			// Verify request configuration
			if req.Method != tt.method {
				t.Errorf("Expected method %s, got %s", tt.method, req.Method)
			}
		})
	}
}
