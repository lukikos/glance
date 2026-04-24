package glance

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/glanceapp/glance/internal/config"
)

// Application holds the core application state and dependencies.
type Application struct {
	Config *config.Config
	server *http.Server
}

// New creates a new Application instance from the provided configuration.
func New(cfg *config.Config) (*Application, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config must not be nil")
	}

	app := &Application{
		Config: cfg,
	}

	return app, nil
}

// Start initializes the HTTP server and begins serving requests.
// It blocks until the context is cancelled or an error occurs.
func (a *Application) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	a.registerRoutes(mux)

	addr := fmt.Sprintf("%s:%d", a.Config.Server.Host, a.Config.Server.Port)

	a.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErr := make(chan error, 1)

	go func() {
		slog.Info("Starting server", "address", addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		return a.Shutdown()
	}
}

// Shutdown gracefully stops the HTTP server with a timeout.
func (a *Application) Shutdown() error {
	if a.server == nil {
		return nil
	}

	slog.Info("Shutting down server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	slog.Info("Server stopped")
	return nil
}

// registerRoutes sets up all HTTP route handlers on the provided mux.
func (a *Application) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", a.handleHealth)
}

// handleHealth responds with a simple health check payload.
func (a *Application) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
