package weather

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type currentReader interface {
	Current(ctx context.Context, zipcode string) (Temperatures, error)
}

// Handler serves GET /weather/{zipcode}.
type Handler struct {
	service currentReader
	logger  *slog.Logger
}

// NewHandler returns the HTTP handler for the weather endpoint.
func NewHandler(service currentReader, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type temperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// ServeHTTP validates the CEP and writes the temperature contract.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zipcode := r.PathValue("zipcode")
	temps, err := h.service.Current(r.Context(), zipcode)
	if err != nil {
		h.writeError(w, zipcode, err)
		return
	}

	writeJSON(w, http.StatusOK, temperatureResponse{
		TempC: temps.Celsius,
		TempF: temps.Fahrenheit,
		TempK: temps.Kelvin,
	})
}

func (h *Handler) writeError(w http.ResponseWriter, zipcode string, err error) {
	switch {
	case errors.Is(err, ErrInvalidZipcode):
		writeText(w, http.StatusUnprocessableEntity, ErrInvalidZipcode.Error())
	case errors.Is(err, ErrZipcodeNotFound):
		writeText(w, http.StatusNotFound, ErrZipcodeNotFound.Error())
	default:
		h.logger.Error("weather request failed", "zipcode", zipcode, "error", err)
		writeText(w, http.StatusBadGateway, "weather service unavailable")
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		writeText(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func writeText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
}
