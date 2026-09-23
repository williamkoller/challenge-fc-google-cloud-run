package weatherapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/williamkoller/challenge-fc-google-cloud-run/internal/weather"
)

func TestClientCurrentCelsius(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		status  int
		body    string
		want    float64
		wantErr error
	}{
		{
			name:   "success",
			status: http.StatusOK,
			body:   `{"current":{"temp_c":28.5}}`,
			want:   28.5,
		},
		{
			name:    "missing temperature",
			status:  http.StatusOK,
			body:    `{"current":{}}`,
			wantErr: weather.ErrUpstream,
		},
		{
			name:    "provider error",
			status:  http.StatusBadRequest,
			body:    `{"error":{"message":"No matching location"}}`,
			wantErr: weather.ErrUpstream,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/current.json" {
					t.Errorf("path = %s", r.URL.Path)
				}
				query := r.URL.Query()
				if query.Get("key") != "test-key" {
					t.Errorf("key was not forwarded to the weather host")
				}
				if query.Get("q") != "São Paulo, SP, Brazil" {
					t.Errorf("q = %q", query.Get("q"))
				}
				if query.Get("aqi") != "no" {
					t.Errorf("aqi = %q", query.Get("aqi"))
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(server.Close)

			client, err := New(server.URL, "test-key")
			if err != nil {
				t.Fatalf("New: %v", err)
			}

			got, err := client.CurrentCelsius(context.Background(), "São Paulo, SP, Brazil")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("celsius = %v, want %v", got, tt.want)
			}
			if err != nil && strings.Contains(err.Error(), "test-key") {
				t.Fatalf("error leaked the api key: %v", err)
			}
		})
	}
}

func TestNewRequiresAPIKey(t *testing.T) {
	t.Parallel()

	if _, err := New("https://api.weatherapi.com", " "); err == nil {
		t.Fatal("expected missing api key to fail")
	}
}
