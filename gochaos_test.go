package gochaos

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew_InvalidTargetURL(t *testing.T) {
	cfg := DefaultConfig()
	cfg.TargetURL = "://bad url"

	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error for invalid target_url, got nil")
	}
}

func TestHandler_ProxiesToTarget(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("upstream-ok"))
	}))
	defer upstream.Close()

	cfg := DefaultConfig()
	cfg.TargetURL = upstream.URL
	cfg.Chaos.ErrorRate = 0
	cfg.Chaos.DisconnectRate = 0
	cfg.Chaos.DNSFailureRate = 0
	cfg.Chaos.UpstreamTimeoutRate = 0
	cfg.Chaos.LatencyMs = 0
	cfg.Chaos.LatencyMinMs = 0
	cfg.Chaos.LatencyMaxMs = 0

	app, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	proxy := httptest.NewServer(app.Handler())
	defer proxy.Close()

	resp, err := http.Get(proxy.URL + "/api/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}
}
