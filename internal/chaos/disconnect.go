package chaos

import (
	"math/rand"
	"net/http"

	"go-chaos/internal/config"
)

func MaybeDisconnect(cfg config.Config, w http.ResponseWriter) bool {
	if rand.Float64() < cfg.Chaos.DisconnectRate {
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, err := hj.Hijack()
			if err == nil {
				_ = conn.Close()
				return true
			}
		}

		w.WriteHeader(http.StatusServiceUnavailable)
		return true
	}
	return false
}
