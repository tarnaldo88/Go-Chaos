package chaos

import (
	"net/http"
	"strings"

	"go-chaos/internal/config"
	"go-chaos/internal/observability"
)

type Store interface {
	Get() config.Config
}

func Middleware(store Store, log *observability.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := store.Get()
		if !ShouldApply(cfg, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if cfg.Chaos.LatencyMs > 0 {
			if log != nil {
				log.Printf("chaos latency=%dms path=%s", cfg.Chaos.LatencyMs, r.URL.Path)
			}
		}

		// latency
		MaybeSleep(cfg)

		// error response
		if MaybeReturnError(cfg, w) {
			if log != nil {
				log.Printf("chaos error status=503 path=%s", r.URL.Path)
			}
			return
		}

		// disconnect
		if MaybeDisconnect(cfg, w) {
			if log != nil {
				log.Printf("chaos disconnect path=%s", r.URL.Path)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ShouldApply(cfg config.Config, path string) bool {
	if path == "" {
		return true
	}
	for _, p := range cfg.Chaos.ExcludePaths {
		if strings.HasPrefix(path, p) {
			return false
		}
	}
	if len(cfg.Chaos.IncludePaths) == 0 {
		return true
	}
	for _, p := range cfg.Chaos.IncludePaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
