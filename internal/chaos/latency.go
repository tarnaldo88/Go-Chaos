package chaos

import (
	"math/rand"
	"time"

	"go-chaos/internal/config"
)

func MaybeSleep(cfg config.Config) (time.Duration, bool) {
	minMs := cfg.Chaos.LatencyMinMs
	maxMs := cfg.Chaos.LatencyMaxMs

	if minMs > 0 || maxMs > 0 {
		if maxMs < minMs {
			return 0, false
		}
		delayMs := minMs
		if maxMs > minMs {
			delayMs = minMs + rand.Intn(maxMs-minMs+1)
		}
		d := time.Duration(delayMs) * time.Millisecond
		time.Sleep(d)
		return d, true
	}

	if cfg.Chaos.LatencyMs <= 0 {
		return 0, false
	}
	d := time.Duration(cfg.Chaos.LatencyMs) * time.Millisecond
	time.Sleep(d)
	return d, true
}
