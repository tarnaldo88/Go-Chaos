package config

import "os"

// import(
// 	"errors"
// 	"io"
// 	"os"
// 	"strings"
// 	"sync/atomic"
// 	"gopkg.in/yaml.v3"
// )

type Config struct {
	ListenAddr string      `yaml:"listen_addr"`
	TargetURL  string      `yaml:"target_url"`
	Chaos      ChaosConfig `yaml:"chaos"`
}

type ChaosConfig struct {
	ErrorRate      float64 `yaml:"error_rate"`
	DisconnectRate float64 `yaml:"disconnect_rate"`
	LatencyMs      int     `yaml:"latency_ms"`
}

func Default() Config {
	return Config{
		ListenAddr: ":8080",
		TargetURL:  "http//localhost:9000",
		Chaos: ChaosConfig{
			ErrorRate:      0.0,
			DisconnectRate: 0.0,
			LatencyMs:      0,
		},
	}
}

func LoadFromFile(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()
	return decodeYAML(f)
}
