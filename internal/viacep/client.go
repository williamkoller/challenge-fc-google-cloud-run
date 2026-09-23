package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/williamkoller/challenge-fc-google-cloud-run/internal/platform/httpclient"
	"github.com/williamkoller/challenge-fc-google-cloud-run/internal/weather"
)

const maxBodyBytes = 1 << 20

// Client looks up a Brazilian CEP on ViaCEP.
type Client struct {
	baseURL string
	http    *http.Client
}

// New builds a ViaCEP client. baseURL defaults to https://viacep.com.br when empty.
func New(baseURL string) (*Client, error) {
	if baseURL == "" {
		baseURL = "https://viacep.com.br"
	}
	baseURL = strings.TrimRight(baseURL, "/")
	httpClient, parsed, err := httpclient.New(baseURL)
	if err != nil {
		return nil, fmt.Errorf("viacep client: %w", err)
	}
	return &Client{
		baseURL: parsed.String(),
		http:    httpClient,
	}, nil
}

type addressPayload struct {
	Localidade string          `json:"localidade"`
	UF         string          `json:"uf"`
	Erro       json.RawMessage `json:"erro"`
}

// Locate returns the city for zipcode.
func (c *Client) Locate(ctx context.Context, zipcode string) (weather.Location, error) {
	if err := weather.ValidateZipcode(zipcode); err != nil {
		return weather.Location{}, err
	}

	endpoint := c.baseURL + "/ws/" + zipcode + "/json/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return weather.Location{}, fmt.Errorf("building viacep request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return weather.Location{}, fmt.Errorf("calling viacep: %w", weather.ErrUpstream)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return weather.Location{}, fmt.Errorf("reading viacep response: %w", weather.ErrUpstream)
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
		return weather.Location{}, weather.ErrInvalidZipcode
	case http.StatusNotFound:
		return weather.Location{}, weather.ErrZipcodeNotFound
	default:
		return weather.Location{}, fmt.Errorf("viacep status %d: %w", resp.StatusCode, weather.ErrUpstream)
	}

	var payload addressPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return weather.Location{}, fmt.Errorf("decoding viacep response: %w", weather.ErrUpstream)
	}
	if erroTrue(payload.Erro) || strings.TrimSpace(payload.Localidade) == "" {
		return weather.Location{}, weather.ErrZipcodeNotFound
	}

	return weather.Location{
		City:  payload.Localidade,
		State: payload.UF,
	}, nil
}

var _ weather.Locator = (*Client)(nil)

func erroTrue(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var asBool bool
	if err := json.Unmarshal(raw, &asBool); err == nil {
		return asBool
	}
	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		return strings.EqualFold(asString, "true")
	}
	return false
}
