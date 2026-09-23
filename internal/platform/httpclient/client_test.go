package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRejectsRedirectToAnotherHost(t *testing.T) {
	t.Parallel()

	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(other.Close)

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL+"/secret", http.StatusFound)
	}))
	t.Cleanup(origin.Close)

	client, _, err := New(origin.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	resp, err := client.Get(origin.URL)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("status = %d, want redirect refused", resp.StatusCode)
	}
	if resp.Header.Get("Location") == "" {
		t.Fatal("expected the original redirect to be returned without following it")
	}
}
