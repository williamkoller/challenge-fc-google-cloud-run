package weather

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Locator resolves a CEP to a city.
type Locator interface {
	Locate(ctx context.Context, zipcode string) (Location, error)
}

// Thermometer reads the current temperature in Celsius for a place query.
type Thermometer interface {
	CurrentCelsius(ctx context.Context, query string) (float64, error)
}

// Service loads the city for a CEP and converts the current temperature.
type Service struct {
	locator     Locator
	thermometer Thermometer
}

// NewService wires the location and weather providers.
func NewService(locator Locator, thermometer Thermometer) *Service {
	return &Service{
		locator:     locator,
		thermometer: thermometer,
	}
}

// Current returns the temperature for an 8-digit CEP.
func (s *Service) Current(ctx context.Context, zipcode string) (Temperatures, error) {
	if err := ValidateZipcode(zipcode); err != nil {
		return Temperatures{}, err
	}

	location, err := s.locator.Locate(ctx, zipcode)
	if err != nil {
		if errors.Is(err, ErrInvalidZipcode) || errors.Is(err, ErrZipcodeNotFound) {
			return Temperatures{}, err
		}
		return Temperatures{}, fmt.Errorf("locating zipcode: %w", err)
	}
	if strings.TrimSpace(location.City) == "" {
		return Temperatures{}, ErrZipcodeNotFound
	}

	celsius, err := s.thermometer.CurrentCelsius(ctx, location.Query())
	if err != nil {
		return Temperatures{}, fmt.Errorf("reading temperature: %w", err)
	}
	return FromCelsius(celsius), nil
}
