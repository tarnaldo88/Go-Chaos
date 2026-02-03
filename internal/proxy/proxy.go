package proxy

import (
	"net/http/httputil"
	"net/url"

	"go-chaos/internal/config"
)

type Store interface {
	Get() config.Config
}

func NewReverseProxy(store Store) (*httputil.ReverseProxy, error) {
	cfg := store.Get()
	target, err := url.Parse(cfg.TargetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Transport with sane defaults.
	base := NewTransport()
	proxy.Transport = NewChaosRoundTripper(base, store)

	return proxy, nil
}
