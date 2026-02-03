package server

import (
	"io"
	"net/http"

	"go-chaos/internal/chaos"
	"go-chaos/internal/config"
	"go-chaos/internal/observability"
	"go-chaos/internal/proxy"

	"gopkg.in/yaml.v3"
)

type Server struct {
	cfg   *config.Store
	log   *observability.Logger
	mux   *http.ServeMux
	proxy http.Handler
}

func New(cfg *config.Store, log *observability.Logger) *Server {
	s := &Server{
		cfg: cfg,
		log: log,
		mux: http.NewServeMux(),
	}

	p, err := proxy.NewReverseProxy(cfg.Get())

	if err != nil {
		panic(err)
	}
	s.proxy = chaos.Middleware(cfg, p)

	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/admin/config", s.handleConfig)
	s.mux.HandleFunc("/healthz", s.handleHealth)

	s.mux.Handle("/", s.proxy)
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleConfigUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cfg, err := config.LoadFromBytes(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if err := s.cfg.Set(cfg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("updated"))
}

func (s *Server) handleConfigGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cfg := s.cfg.Get()

	out, err := yaml.Marshal(cfg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out)
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleConfigGet(w, r)
	case http.MethodPost:
		s.handleConfigUpdate(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
