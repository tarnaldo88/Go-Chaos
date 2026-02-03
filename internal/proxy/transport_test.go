package proxy

import (
	"net"
	"net/http"
	"net/url"
	"testing"

	"go-chaos/internal/config"
)

type staticStore struct {
	cfg config.Config
}

func (s staticStore) Get() config.Config { return s.cfg }

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestChaosRoundTripper_DNSFailure(t *testing.T) {
	cfg := config.Default()
	cfg.Chaos.DNSFailureRate = 1.0
	cfg.Chaos.IncludePaths = []string{"/api/"}

	store := staticStore{cfg: cfg}

	base := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		t.Fatalf("base RoundTripper should not be called")
		return nil, nil
	})

	rt := NewChaosRoundTripper(base, store)

	u, _ := url.Parse("http://example.com/api/test")
	req := &http.Request{URL: u}

	_, err := rt.RoundTrip(req)
	if err == nil {
		t.Fatalf("expected DNS error, got nil")
	}
	if _, ok := err.(*net.DNSError); !ok {
		t.Fatalf("expected *net.DNSError, got %T", err)
	}
}

func TestChaosRoundTripper_Timeout(t *testing.T) {
	cfg := config.Default()
	cfg.Chaos.UpstreamTimeoutRate = 1.0
	cfg.Chaos.IncludePaths = []string{"/api/"}

	store := staticStore{cfg: cfg}

	base := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		t.Fatalf("base RoundTripper should not be called")
		return nil, nil
	})

	rt := NewChaosRoundTripper(base, store)

	u, _ := url.Parse("http://example.com/api/test")
	req := &http.Request{URL: u}

	_, err := rt.RoundTrip(req)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}

	netErr, ok := err.(net.Error)
	if !ok || !netErr.Timeout() {
		t.Fatalf("expected timeout net.Error, got %T", err)
	}
}
