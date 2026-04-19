package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/glanceapp/glance/internal/config"
	"github.com/glanceapp/glance/internal/server"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	configPath := flag.String("config", "glance.yml", "Path to the configuration file")
	showVersion := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("glance %s (commit: %s, built: %s)\n", Version, Commit, Date)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	log.Printf("Starting glance %s on %s:%d", Version, cfg.Server.Host, cfg.Server.Port)

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
