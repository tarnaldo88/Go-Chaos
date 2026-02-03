package chaos

import (
	"time"

	"go-chaos/internal/config"
)

func MaybeSleep(cfg config.Config) {
	if cfg.Chaos.LatencyMs <= 0 {
		return
	}
	time.Sleep(time.Duration(cfg.Chaos.LatencyMs) * time.Millisecond)
}
