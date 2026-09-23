package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Config is the process configuration loaded from the environment.
type Config struct {
	Port           string
	WeatherAPIKey  string
	ViaCEPBaseURL  string
	WeatherBaseURL string
}

// Load reads configuration from the environment.
// WEATHER_API_KEY is required. PORT defaults to 8080.
func Load() (Config, error) {
	apiKey := strings.TrimSpace(os.Getenv("WEATHER_API_KEY"))
	if apiKey == "" {
		return Config{}, errors.New("WEATHER_API_KEY is required")
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	viaCEP := strings.TrimRight(envOr("VIACEP_BASE_URL", "https://viacep.com.br"), "/")
	weatherAPI := strings.TrimRight(envOr("WEATHER_BASE_URL", "https://api.weatherapi.com"), "/")
	if err := validateBaseURL(viaCEP); err != nil {
		return Config{}, fmt.Errorf("VIACEP_BASE_URL: %w", err)
	}
	if err := validateBaseURL(weatherAPI); err != nil {
		return Config{}, fmt.Errorf("WEATHER_BASE_URL: %w", err)
	}

	return Config{
		Port:           port,
		WeatherAPIKey:  apiKey,
		ViaCEPBaseURL:  viaCEP,
		WeatherBaseURL: weatherAPI,
	}, nil
}

func envOr(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func validateBaseURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing url: %w", err)
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return errors.New("scheme must be http or https")
	}
	if parsed.Host == "" {
		return errors.New("host is required")
	}
	return nil
}
