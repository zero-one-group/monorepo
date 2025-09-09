package server

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"go-modular/internal/adapter"
)

// HTTPServer is the main HTTP server struct.
// Logger and Tracer are injected for observability.
type HTTPServer struct {
	httpAddr string
	logger   *slog.Logger
}

func NewHTTPServer(httpAddr string) *HTTPServer {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	return &HTTPServer{
		httpAddr: httpAddr,
		logger:   logger,
	}
}

func (s *HTTPServer) Start() error {
	databaseURL := os.Getenv("DATABASE_URL")

	// Initialize Postgres pool
	pg, err := adapter.NewPostgres(adapter.PostgresConfig{
		URL:        databaseURL,
		SearchPath: "public",
	})
	if err != nil {
		slog.Error("Failed to connect to Postgres database", "err", err)
		os.Exit(1)
	}
	defer pg.Close()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	h := NewServerHandler(pg.Pool, s.logger)
	h.RegisterRoutes(e)

	s.logger.Info("Starting HTTP server", "addr", s.httpAddr)

	// Start the server
	return e.Start(s.httpAddr)
}
