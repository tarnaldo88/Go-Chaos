package proxy

import (
	"net/http/httputil"
	"net/url"

	"go-chaos/internal/config"
)

func NewReverseProxy(cfg config.Config) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(cfg.TargetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Transport with sane defaults.
	proxy.Transport = NewTransport()

	return proxy, nil
}
