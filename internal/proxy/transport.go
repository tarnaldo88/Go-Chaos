package proxy

import (
	"math/rand"
	"net"
	"net/http"
	"time"

	"go-chaos/internal/config"
)

type Store interface {
	Get() config.Config
}

func NewTransport() *http.Transport {
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

type ChaosRoundTripper struct {
	base  http.RoundTripper
	store Store
}

func NewChaosRoundTripper(base http.RoundTripper, store Store) *ChaosRoundTripper {
	return &ChaosRoundTripper{base: base, store: store}
}

func (c *ChaosRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	cfg := c.store.Get()

	if cfg.Chaos.DNSFailureRate > 0 && rand.Float64() < cfg.Chaos.DNSFailureRate {
		return nil, &net.DNSError{
			Err:        "no such host",
			Name:       r.URL.Hostname(),
			IsNotFound: true,
		}
	}

	if cfg.Chaos.UpstreamTimeoutRate > 0 && rand.Float64() < cfg.Chaos.UpstreamTimeoutRate {
		return nil, timeoutError("upstream timeout (chaos)")
	}

	return c.base.RoundTrip(r)
}

type timeoutError string

func (e timeoutError) Error() string   { return string(e) }
func (e timeoutError) Timeout() bool   { return true }
func (e timeoutError) Temporary() bool { return true }
