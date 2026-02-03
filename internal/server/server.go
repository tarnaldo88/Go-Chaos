package server

import (
	"io"
	"net/http"

	"go-chaos/internal/chaos"
	"go-chaos/internal/config"
	"go-chaos/internal/observability"
	"go-chaos/internal/proxy"
)

type Server struct {
	cfg *config.Store
	log *observability.Logger
	mux *http.ServeMux
}

func New(cfg *config.Store, log *observability.Logger) Server {
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
	s.mux.HandleFunc("/admin/config", s.handleConfigUpdate)
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
