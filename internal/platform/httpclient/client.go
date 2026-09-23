package httpclient

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// New returns an HTTP client that only follows redirects to host.
func New(rawBaseURL string) (*http.Client, *url.URL, error) {
	parsed, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing base url: %w", err)
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return nil, nil, errors.New("base url scheme must be http or https")
	}
	if parsed.Host == "" {
		return nil, nil, errors.New("base url host is required")
	}

	host := parsed.Host
	client := &http.Client{
		Timeout: 4 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return errors.New("stopped after 3 redirects")
			}
			if req.URL.Host != host {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	return client, parsed, nil
}
