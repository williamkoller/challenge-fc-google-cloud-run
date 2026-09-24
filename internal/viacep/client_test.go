package viacep

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/williamkoller/challenge-fc-google-cloud-run/internal/weather"
)

func TestClientLocate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		zipcode    string
		status     int
		body       string
		want       weather.Location
		wantErr    error
		wantCalled bool
	}{
		{
			name:       "found",
			zipcode:    "01001000",
			status:     http.StatusOK,
			body:       `{"localidade":"São Paulo","uf":"SP"}`,
			want:       weather.Location{City: "São Paulo", State: "SP"},
			wantCalled: true,
		},
		{
			name:       "boolean error flag",
			zipcode:    "99999999",
			status:     http.StatusOK,
			body:       `{"erro":true}`,
			wantErr:    weather.ErrZipcodeNotFound,
			wantCalled: true,
		},
		{
			name:       "string error flag",
			zipcode:    "99999999",
			status:     http.StatusOK,
			body:       `{"erro":"true"}`,
			wantErr:    weather.ErrZipcodeNotFound,
			wantCalled: true,
		},
		{
			name:       "provider rejects format",
			zipcode:    "01001000",
			status:     http.StatusBadRequest,
			body:       `{}`,
			wantErr:    weather.ErrInvalidZipcode,
			wantCalled: true,
		},
		{
			name:       "provider bad gateway is unknown zipcode",
			zipcode:    "00000001",
			status:     http.StatusBadGateway,
			body:       `<html><head><title>502 Bad Gateway</title></head></html>`,
			wantErr:    weather.ErrZipcodeNotFound,
			wantCalled: true,
		},
		{
			name:       "provider unavailable",
			zipcode:    "01001000",
			status:     http.StatusInternalServerError,
			body:       `{}`,
			wantErr:    weather.ErrUpstream,
			wantCalled: true,
		},
		{
			name:    "invalid zipcode never calls provider",
			zipcode: "123",
			wantErr: weather.ErrInvalidZipcode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var called atomic.Bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called.Store(true)
				if r.URL.Path != "/ws/"+tt.zipcode+"/json/" {
					t.Errorf("path = %s", r.URL.Path)
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(server.Close)

			client, err := New(server.URL)
			if err != nil {
				t.Fatalf("New: %v", err)
			}

			got, err := client.Locate(context.Background(), tt.zipcode)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if called.Load() != tt.wantCalled {
				t.Fatalf("called = %v, want %v", called.Load(), tt.wantCalled)
			}
			if err == nil && got != tt.want {
				t.Fatalf("location = %+v, want %+v", got, tt.want)
			}
		})
	}
}
