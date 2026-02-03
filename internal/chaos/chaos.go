package chaos

import (
	"net/http"

	"go-chaos/internal/config"
	"go-chaos/internal/observability"
)

type Store interface {
	Get() config.Config
}

func Middleware(store Store, log *observability.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := store.Get()

		if cfg.Chaos.LatencyMs > 0 {
			log.Printf("chaos latency=%dms path=%s", cfg.Chaos.LatencyMs, r.URL.Path)
		}

		// latency
		MaybeSleep(cfg)

		// error response
		if MaybeReturnError(cfg, w) {
			log.Printf("chaos error status=503 path=%s", r.URL.Path)
			return
		}

		// disconnect
		if MaybeDisconnect(cfg, w) {
			log.Printf("chaos disconnect path=%s", r.URL.Path)
			return
		}

		next.ServeHTTP(w, r)
	})
}
