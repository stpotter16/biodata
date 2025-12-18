package parse

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseLoginPost(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		wantErr     bool
		wantPass    string
	}{
		{
			name:        "valid passphrase",
			requestBody: `{"passphrase": "my-secret-pass"}`,
			wantErr:     false,
			wantPass:    "my-secret-pass",
		},
		{
			name:        "empty passphrase",
			requestBody: `{"passphrase": ""}`,
			wantErr:     false,
			wantPass:    "",
		},
		{
			name:        "malformed JSON",
			requestBody: `{"passphrase": "test"`,
			wantErr:     true,
		},
		{
			name:        "missing passphrase field",
			requestBody: `{"username": "test"}`,
			wantErr:     false,
			wantPass:    "",
		},
		{
			name:        "empty request body",
			requestBody: ``,
			wantErr:     true,
		},
		{
			name:        "null passphrase",
			requestBody: `{"passphrase": null}`,
			wantErr:     false,
			wantPass:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			got, err := ParseLoginPost(req)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLoginPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.Passphrase != tt.wantPass {
				t.Errorf("ParseLoginPost() passphrase = %v, want %v", got.Passphrase, tt.wantPass)
			}
		})
	}
}
