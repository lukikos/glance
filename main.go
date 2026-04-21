package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/glanceapp/glance/internal/config"
	"github.com/glanceapp/glance/internal/server"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	// Default config path changed to ~/.config/glance/glance.yml for XDG compliance
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	defaultConfigPath := filepath.Join(homeDir, ".config", "glance", "glance.yml")

	configPath := flag.String("config", defaultConfigPath, "Path to the configuration file")
	showVersion := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("glance %s (commit: %s, built: %s)\n", Version, Commit, Date)
		os.Exit(0)
	}

	// Warn the user if the config file doesn't exist at the resolved path,
	// so they get a clear message instead of a cryptic load error.
	if _, statErr := os.Stat(*configPath); os.IsNotExist(statErr) {
		log.Printf("Warning: config file not found at %s", *configPath)
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
