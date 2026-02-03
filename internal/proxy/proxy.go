package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"go-chaos/internal/config"
)

func NewReverseProxy(cfg config.Config) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(cfg.TargetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Transport with sane defaults.
	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return proxy, nil
}
