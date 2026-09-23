package config

import "testing"

func TestLoadRequiresAPIKey(t *testing.T) {
	t.Setenv("WEATHER_API_KEY", "")
	t.Setenv("PORT", "")
	t.Setenv("VIACEP_BASE_URL", "")
	t.Setenv("WEATHER_BASE_URL", "")

	if _, err := Load(); err == nil {
		t.Fatal("expected missing WEATHER_API_KEY to fail")
	}
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("WEATHER_API_KEY", "secret")
	t.Setenv("PORT", "")
	t.Setenv("VIACEP_BASE_URL", "")
	t.Setenv("WEATHER_BASE_URL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Port != "8080" {
		t.Fatalf("port = %q", cfg.Port)
	}
	if cfg.ViaCEPBaseURL != "https://viacep.com.br" {
		t.Fatalf("viacep = %q", cfg.ViaCEPBaseURL)
	}
	if cfg.WeatherBaseURL != "https://api.weatherapi.com" {
		t.Fatalf("weather = %q", cfg.WeatherBaseURL)
	}
	if cfg.WeatherAPIKey != "secret" {
		t.Fatal("api key was not loaded")
	}
}

func TestLoadRejectsBadBaseURL(t *testing.T) {
	t.Setenv("WEATHER_API_KEY", "secret")
	t.Setenv("VIACEP_BASE_URL", "ftp://example.com")

	if _, err := Load(); err == nil {
		t.Fatal("expected invalid base url to fail")
	}
}
