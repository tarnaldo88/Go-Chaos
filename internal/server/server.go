package server

// import (
// 	"io"
// 	"net/http"

// 	"go-chaos/internal/config"
// 	"go-chaos/internal/observability"
// )

type Server struct {
	cfg *config.Store
	log *observability.Logger
	mux *http.ServeMux
}

func New(cfg *config.Store, log *observability.Logger) Server {
	s := &Server{
		cfg: cfg,
		log: log,
		mux: http.NewServerMux(),
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/admin/config", s.handleConfigUpdate)
	s.mux.HandleFunc("/healthz", s.handleHealth)
	// TODO: add proxy handler (wrap with chaos middleware)
}

func (s *Server) Handler() http.Handler {
	return s.mux
}
