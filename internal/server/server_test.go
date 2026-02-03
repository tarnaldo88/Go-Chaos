package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-chaos/internal/config"
	"go-chaos/internal/observability"
)

func newTestServer(t *testing.T, cfg config.Config) *httptest.Server {
	t.Helper()
	store := config.NewStore(cfg)
	log := observability.New()

	s := New(store, log)
	return httptest.NewServer(s.Handler())
}

func TestLatencyInjection(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	cfg := config.Default()
	cfg.TargetURL = upstream.URL
	cfg.Chaos.LatencyMs = 150

	proxy := newTestServer(t, cfg)
	defer proxy.Close()

	start := time.Now()
	resp, err := http.Get(proxy.URL + "/ping")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	_ = resp.Body.Close()

	elapsed := time.Since(start)
	if elapsed < 140*time.Millisecond {
		t.Fatalf("expected latency >= 140ms, got %v", elapsed)
	}
}

func TestErrorInjection(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	cfg := config.Default()
	cfg.TargetURL = upstream.URL
	cfg.Chaos.ErrorRate = 1.0 // always error

	proxy := newTestServer(t, cfg)
	defer proxy.Close()

	resp, err := http.Get(proxy.URL + "/ping")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", resp.StatusCode)
	}
}

func TestConfigUpdateEndpoint(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	cfg := config.Default()
	cfg.TargetURL = upstream.URL

	proxy := newTestServer(t, cfg)
	defer proxy.Close()

	updated := []byte(`
listen_addr: ":8080"
target_url: "` + upstream.URL + `"
chaos:
  error_rate: 0.5
  disconnect_rate: 0.0
  latency_ms: 0
`)
	resp, err := http.Post(proxy.URL+"/admin/config", "application/x-yaml", bytes.NewReader(updated))
	if err != nil {
		t.Fatalf("config update failed: %v", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
