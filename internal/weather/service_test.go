package weather

import (
	"context"
	"errors"
	"testing"
)

type stubLocator struct {
	location Location
	err      error
	called   bool
}

func (s *stubLocator) Locate(context.Context, string) (Location, error) {
	s.called = true
	return s.location, s.err
}

type stubThermometer struct {
	celsius float64
	err     error
	query   string
}

func (s *stubThermometer) CurrentCelsius(_ context.Context, query string) (float64, error) {
	s.query = query
	return s.celsius, s.err
}

func TestServiceCurrent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		zipcode     string
		locator     *stubLocator
		thermometer *stubThermometer
		want        Temperatures
		wantErr     error
		wantCalled  bool
		wantQuery   string
	}{
		{
			name:    "success",
			zipcode: "01001000",
			locator: &stubLocator{location: Location{City: "São Paulo", State: "SP"}},
			thermometer: &stubThermometer{
				celsius: 28.5,
			},
			want:       FromCelsius(28.5),
			wantCalled: true,
			wantQuery:  "São Paulo, SP, Brazil",
		},
		{
			name:        "invalid zipcode skips providers",
			zipcode:     "123",
			locator:     &stubLocator{},
			thermometer: &stubThermometer{},
			wantErr:     ErrInvalidZipcode,
		},
		{
			name:        "unknown zipcode",
			zipcode:     "00000000",
			locator:     &stubLocator{err: ErrZipcodeNotFound},
			thermometer: &stubThermometer{},
			wantErr:     ErrZipcodeNotFound,
			wantCalled:  true,
		},
		{
			name:        "empty city",
			zipcode:     "01001000",
			locator:     &stubLocator{location: Location{}},
			thermometer: &stubThermometer{},
			wantErr:     ErrZipcodeNotFound,
			wantCalled:  true,
		},
		{
			name:    "weather provider failure",
			zipcode: "01001000",
			locator: &stubLocator{location: Location{City: "São Paulo", State: "SP"}},
			thermometer: &stubThermometer{
				err: ErrUpstream,
			},
			wantErr:    ErrUpstream,
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewService(tt.locator, tt.thermometer)
			got, err := svc.Current(context.Background(), tt.zipcode)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if tt.locator.called != tt.wantCalled {
				t.Fatalf("locator called = %v, want %v", tt.locator.called, tt.wantCalled)
			}
			if err != nil {
				return
			}
			if got != tt.want {
				t.Fatalf("temperatures = %+v, want %+v", got, tt.want)
			}
			if tt.thermometer.query != tt.wantQuery {
				t.Fatalf("query = %q, want %q", tt.thermometer.query, tt.wantQuery)
			}
		})
	}
}
