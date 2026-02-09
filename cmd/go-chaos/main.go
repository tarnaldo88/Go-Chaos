package main

import (
	"flag"
	"log"
	gochaos "go-chaos"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := gochaos.LoadConfigFile(*configPath)
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	app, err := gochaos.New(cfg)
	if err != nil {
		log.Fatalf("server setup failed: %v", err)
	}

	if err := app.ListenAndServe(""); err != nil {
		log.Fatal(err)
	}
}
