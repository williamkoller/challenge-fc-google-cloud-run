package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/williamkoller/challenge-fc-google-cloud-run/internal/platform/httpclient"
	"github.com/williamkoller/challenge-fc-google-cloud-run/internal/weather"
)

const maxBodyBytes = 1 << 20

// Client reads the current temperature from WeatherAPI.
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// New builds a WeatherAPI client. baseURL defaults to https://api.weatherapi.com when empty.
func New(baseURL, apiKey string) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("weatherapi client: api key is required")
	}
	if baseURL == "" {
		baseURL = "https://api.weatherapi.com"
	}
	baseURL = strings.TrimRight(baseURL, "/")
	httpClient, parsed, err := httpclient.New(baseURL)
	if err != nil {
		return nil, fmt.Errorf("weatherapi client: %w", err)
	}
	return &Client{
		baseURL: parsed.String(),
		apiKey:  apiKey,
		http:    httpClient,
	}, nil
}

type currentPayload struct {
	Current *struct {
		TempC *float64 `json:"temp_c"`
	} `json:"current"`
}

// CurrentCelsius returns the current temperature for query.
// The API key is never included in returned errors.
func (c *Client) CurrentCelsius(ctx context.Context, query string) (float64, error) {
	endpoint, err := url.Parse(c.baseURL + "/v1/current.json")
	if err != nil {
		return 0, fmt.Errorf("building weatherapi url: %w", err)
	}
	values := endpoint.Query()
	values.Set("key", c.apiKey)
	values.Set("q", query)
	values.Set("aqi", "no")
	endpoint.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("building weatherapi request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("calling weatherapi: %w", weather.ErrUpstream)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return 0, fmt.Errorf("reading weatherapi response: %w", weather.ErrUpstream)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi status %d: %w", resp.StatusCode, weather.ErrUpstream)
	}

	var payload currentPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return 0, fmt.Errorf("decoding weatherapi response: %w", weather.ErrUpstream)
	}
	if payload.Current == nil || payload.Current.TempC == nil {
		return 0, fmt.Errorf("weatherapi payload missing temp_c: %w", weather.ErrUpstream)
	}
	return *payload.Current.TempC, nil
}

var _ weather.Thermometer = (*Client)(nil)
