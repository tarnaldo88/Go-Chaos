package chaos

import (
	"math/rand"
	"net/http"

	"go-chaos/internal/config"
)

func MaybeReturnError(cfg config.Config, w http.ResponseWriter) bool {
	if rand.Float64() < cfg.Chaos.ErrorRate {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("chaos: injected error"))
		return true
	}
	return false
}
