package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/williamkoller/challenge-fc-google-cloud-run/internal/platform/config"
	"github.com/williamkoller/challenge-fc-google-cloud-run/internal/platform/httpserver"
	"github.com/williamkoller/challenge-fc-google-cloud-run/internal/viacep"
	"github.com/williamkoller/challenge-fc-google-cloud-run/internal/weather"
	"github.com/williamkoller/challenge-fc-google-cloud-run/internal/weatherapi"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	locator, err := viacep.New(cfg.ViaCEPBaseURL)
	if err != nil {
		return err
	}
	thermometer, err := weatherapi.New(cfg.WeatherBaseURL, cfg.WeatherAPIKey)
	if err != nil {
		return err
	}

	service := weather.NewService(locator, thermometer)
	handler := weather.NewHandler(service, logger)

	mux := http.NewServeMux()
	mux.Handle("GET /weather/{zipcode}", handler)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return httpserver.Run(ctx, ":"+cfg.Port, httpserver.LogRequests(logger, mux), logger)
}
