package config

import (
	"errors"
	"io"
	"os"
	"strings"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

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
		TargetURL:  "http://localhost:9000",
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

func LoadFromBytes(b []byte) (Config, error) {
	return decodeYAML(strings.NewReader(string(b)))
}

func decodeYAML(r io.Reader) (Config, error) {
	cfg := Default()
	dec := yaml.NewDecoder(r)
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, err
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.ListenAddr == "" {
		return errors.New("listen_addr is required")
	}
	if c.TargetURL == "" {
		return errors.New("target_url is required")
	}
	if c.Chaos.ErrorRate < 0 || c.Chaos.ErrorRate > 1 {
		return errors.New("chaos.error_rate must be 0.0-1.0")
	}
	if c.Chaos.DisconnectRate < 0 || c.Chaos.DisconnectRate > 1 {
		return errors.New("chaos.disconnect_rate must be 0.0-1.0")
	}
	if c.Chaos.LatencyMs < 0 {
		return errors.New("chaos.latency_ms must be >= 0")
	}
	return nil
}

type Store struct {
	v atomic.Value
}

func NewStore(cfg Config) *Store {
	s := &Store{}
	s.v.Store(cfg)
	return s
}

func (s *Store) Set(cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	s.v.Store(cfg)
	return nil
}

func (s *Store) Get() Config {
	return s.v.Load().(Config)
}
