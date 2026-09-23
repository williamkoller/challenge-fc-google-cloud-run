package weather

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubService struct {
	temps Temperatures
	err   error
}

func (s stubService) Current(context.Context, string) (Temperatures, error) {
	return s.temps, s.err
}

func TestHandlerServeHTTP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		target     string
		method     string
		service    stubService
		wantStatus int
		wantBody   string
		wantJSON   *temperatureResponse
	}{
		{
			name:   "success",
			target: "/weather/01001000",
			method: http.MethodGet,
			service: stubService{
				temps: FromCelsius(28.5),
			},
			wantStatus: http.StatusOK,
			wantJSON: &temperatureResponse{
				TempC: 28.5,
				TempF: 83.3,
				TempK: 301.5,
			},
		},
		{
			name:       "invalid zipcode",
			target:     "/weather/123",
			method:     http.MethodGet,
			service:    stubService{err: ErrInvalidZipcode},
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   "invalid zipcode",
		},
		{
			name:       "zipcode not found",
			target:     "/weather/00000000",
			method:     http.MethodGet,
			service:    stubService{err: ErrZipcodeNotFound},
			wantStatus: http.StatusNotFound,
			wantBody:   "can not find zipcode",
		},
		{
			name:       "upstream failure",
			target:     "/weather/01001000",
			method:     http.MethodGet,
			service:    stubService{err: errors.New("timeout")},
			wantStatus: http.StatusBadGateway,
			wantBody:   "weather service unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := NewHandler(tt.service, slog.New(slog.DiscardHandler))
			mux := http.NewServeMux()
			mux.Handle("GET /weather/{zipcode}", handler)

			req := httptest.NewRequest(tt.method, tt.target, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d, body = %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantBody != "" && rec.Body.String() != tt.wantBody {
				t.Fatalf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
			if tt.wantJSON == nil {
				return
			}
			var got temperatureResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decoding body: %v", err)
			}
			if got != *tt.wantJSON {
				t.Fatalf("payload = %+v, want %+v", got, *tt.wantJSON)
			}
		})
	}
}

func TestHandlerRejectsOtherMethods(t *testing.T) {
	t.Parallel()

	handler := NewHandler(stubService{}, slog.New(slog.DiscardHandler))
	mux := http.NewServeMux()
	mux.Handle("GET /weather/{zipcode}", handler)

	req := httptest.NewRequest(http.MethodPost, "/weather/01001000", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
