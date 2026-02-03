package chaos

import (
	"net/http"

	"go-chaos/internal/config"
)

type Store interface {
	Get() config.Config
}

func Middleware(store Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := store.Get()

		// 1) latency
		MaybeSleep(cfg)

		// 2) error response
		if MaybeReturnError(cfg, w) {
			return
		}

		// 3) disconnect
		if MaybeDisconnect(cfg, w) {
			return
		}

		next.ServeHTTP(w, r)
	})
}
