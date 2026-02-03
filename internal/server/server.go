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
