package gochaos

import (
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"go-chaos/internal/config"
	"go-chaos/internal/observability"
	"go-chaos/internal/server"
)

type Config = config.Config
type ChaosConfig = config.ChaosConfig

type App struct {
	store *config.Store
	srv   *server.Server
}

var seedOnce sync.Once

func New(cfg Config) (*App, error) {
	seedOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	store := config.NewStore(cfg)
	log := observability.New()

	srv, err := server.New(store, log)
	if err != nil {
		return nil, err
	}

	return &App{
		store: store,
		srv:   srv,
	}, nil
}

func DefaultConfig() Config {
	return config.Default()
}

func LoadConfigFile(path string) (Config, error) {
	return config.LoadFromFile(path)
}

func LoadConfigYAML(b []byte) (Config, error) {
	return config.LoadFromBytes(b)
}

func (a *App) Handler() http.Handler {
	return a.srv.Handler()
}

func (a *App) CurrentConfig() Config {
	return a.store.Get()
}

func (a *App) UpdateConfig(cfg Config) error {
	return a.store.Set(cfg)
}

func (a *App) ListenAndServe(addr string) error {
	if addr == "" {
		addr = a.store.Get().ListenAddr
	}
	log.Printf("go-chaos listening on %s", addr)
	return http.ListenAndServe(addr, a.Handler())
}
