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
		log.Printf("Hint: create the directory with: mkdir -p %s", filepath.Dir(*configPath))
		// Also print an example copy command if the user has a local glance.yml nearby
		if _, localErr := os.Stat("glance.yml"); localErr == nil {
			log.Printf("Hint: a local glance.yml was found; copy it with: cp glance.yml %s", *configPath)
		}
		// Don't fatal here — let config.Load produce the definitive error message
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Log the full address as a clickable URL for convenience.
	// Use 127.0.0.1 instead of localhost to avoid IPv6 resolution delays on
	// some systems where 'localhost' resolves to ::1 first.
	host := cfg.Server.Host
	if host == "" || host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	log.Printf("Starting glance %s — open at http://%s:%d", Version, host, cfg.Server.Port)
	log.Printf("Config loaded from: %s", *configPath)

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
