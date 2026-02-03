package main

import (
	"flag"
	"log"
	"net/http"

	"go-chaos/internal/config"
	"go-chaos/internal/observability"
	"go-chaos/internal/server"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	store := config.NewStore(cfg)
	logger := observability.New()
	srv := server.New(store, logger)

	log.Printf("listening on %s", cfg.ListenAddr)
	if err := http.ListenAndServe(cfg.ListenAddr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}
