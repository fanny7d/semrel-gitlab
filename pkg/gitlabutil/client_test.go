package gitlabutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		apiURL        string
		skipSSLVerify bool
		wantErr       bool
	}{
		{
			name:          "valid token and apiURL",
			token:         "valid-token",
			apiURL:        "https://gitlab.example.com",
			skipSSLVerify: false,
			wantErr:       false,
		},
		{
			name:          "empty token",
			token:         "",
			apiURL:        "https://gitlab.example.com",
			skipSSLVerify: false,
			wantErr:       true,
		},
		{
			name:          "empty apiURL",
			token:         "valid-token",
			apiURL:        "",
			skipSSLVerify: false,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.token, tt.apiURL, tt.skipSSLVerify)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		apiURL        string
		skipSSLVerify bool
		method        string
		path          string
		body          interface{}
		options       []gitlab.RequestOptionFunc
		wantErr       bool
	}{
		{
			name:          "valid request",
			token:         "valid-token",
			apiURL:        "https://gitlab.example.com",
			skipSSLVerify: false,
			method:        "GET",
			path:          "/api/v4/projects",
			body:          nil,
			options:       nil,
			wantErr:       false,
		},
		{
			name:          "invalid URL",
			token:         "valid-token",
			apiURL:        "://invalid-url",
			skipSSLVerify: false,
			method:        "GET",
			path:          "/api/v4/projects",
			body:          nil,
			options:       nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.token, tt.apiURL, tt.skipSSLVerify)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatalf("NewClient() error = %v", err)
			}

			req, err := client.NewRequest(tt.method, tt.path, tt.body, nil)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, req)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
			}
		})
	}
}
