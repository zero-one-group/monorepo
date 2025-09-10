package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"go-modular/docs"
	"go-modular/web"

	"github.com/alexliesenfeld/health"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	scalar "github.com/bdpiprava/scalar-go"
)

// ServerHandler holds dependencies for HTTP handlers.
type ServerHandler struct {
	PGPool *pgxpool.Pool
	Logger *slog.Logger
	WebFS  embed.FS
}

// NewServerHandler creates a new ServerHandler.
func NewServerHandler(pgPool *pgxpool.Pool, logger *slog.Logger) *ServerHandler {
	return &ServerHandler{
		PGPool: pgPool,
		Logger: logger,
		WebFS:  web.WebDir,
	}
}

// RegisterRoutes registers all HTTP routes to the given Echo instance.
func (h *ServerHandler) RegisterRoutes(e *echo.Echo) {
	staticFS := getFileSystem(h.WebFS, "static")
	assetHandler := http.FileServer(staticFS)

	// Serve static files under /static/*
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	// Serve index.html for the root endpoint
	e.GET("/", h.RootHandler(staticFS))
	e.GET("/*", h.RootHandler(staticFS))

	e.GET("/healthz", h.HealthCheckHandler)          // Health check endpoint
	e.GET("/api-docs", h.APIDocsHandler)             // Scalar API docs endpoint
	e.GET("/api/openapi.json", h.OpenAPISpecHandler) // Serve raw OpenAPI spec

}

// RootHandler serves index.html, it can be used for embedded SPA frontend.
func (h *ServerHandler) RootHandler(staticFS http.FileSystem) echo.HandlerFunc {
	return func(c echo.Context) error {
		upath := c.Request().URL.Path
		// Only serve index.html for non-static and non-healthz routes
		if strings.HasPrefix(upath, "/static/") || upath == "/healthz" {
			return echo.ErrNotFound
		}

		f, err := staticFS.Open("index.html")
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "index.html not found")
		}
		defer func() {
			if cerr := f.Close(); cerr != nil {
				h.Logger.Warn("failed to close index.html file", "err", cerr)
			}
		}()

		c.Response().Header().Set(echo.HeaderContentType, "text/html; charset=utf-8")
		http.ServeContent(c.Response(), c.Request(), "index.html", time.Now(), f)

		return nil
	}
}

// @Summary		    Service healthcheck
// @Description	    Checks the health of the service
// @Tags	        General Information
// @Router		    /healthz [get]
func (h *ServerHandler) HealthCheckHandler(c echo.Context) error {
	hc := health.NewChecker(
		health.WithCacheDuration(10*time.Second),
		health.WithTimeout(5*time.Second),
		health.WithCheck(health.Check{
			Name:    "database",
			Timeout: 2 * time.Second,
			Check: func(ctx context.Context) error {
				return h.PGPool.Ping(ctx)
			},
		}),
	)

	// Transform health.NewHandler to Echo handler
	handler := health.NewHandler(hc)
	handler.ServeHTTP(c.Response(), c.Request())

	return nil
}

// OpenAPISpecHandler serves the embedded swagger.json as application/json.
func (h *ServerHandler) OpenAPISpecHandler(c echo.Context) error {
	enableOpenAPI := os.Getenv("APP_ENABLE_API_DOCS")
	if enableOpenAPI != "true" {
		return echo.NewHTTPError(http.StatusNotFound, "API docs are disabled")
	}

	specBytes, err := docs.SwaggerFS.ReadFile("swagger.json")
	if err != nil {
		h.Logger.Error("failed to read swagger.json", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read API spec")
	}
	return c.Blob(http.StatusOK, "application/json", specBytes)
}

// APIDocsHandler serves the Scalar API docs UI using the embedded OpenAPI spec.
func (h *ServerHandler) APIDocsHandler(c echo.Context) error {
	enableOpenAPI := os.Getenv("APP_ENABLE_API_DOCS")
	if enableOpenAPI != "true" {
		return echo.NewHTTPError(http.StatusNotFound, "API docs are disabled")
	}

	// Construct the OpenAPI spec URL based on the current request
	req := c.Request()
	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	specURL := fmt.Sprintf("%s://%s/api/openapi.json", scheme, req.Host)

	// Generate HTML content using scalargo NewV2 with spec URL
	htmlContent, err := scalar.NewV2(
		scalar.WithSpecURL(specURL),
		scalar.WithAuthenticationOpts(
			scalar.WithCustomSecurity(),
			scalar.WithPreferredSecurityScheme("bearerAuth"),
			scalar.WithHTTPBearerToken("your-bearer-token-here"),
		),
		scalar.WithServers(scalar.Server{
			URL:         fmt.Sprintf("%s://%s", scheme, req.Host),
			Description: "Default server",
		}),
		scalar.WithHiddenClients(
			"libcurl",       // C
			"httpclient",    // CSharp
			"restsharp",     // CSharp
			"clj_http",      // Clojure
			"http",          // Dart
			"native",        // Go & Ruby
			"http1.1",       // HTTP
			"asynchttp",     // Java
			"nethttp",       // Java
			"okhttp",        // Java & Kotlin
			"unirest",       // Java
			"jquery",        // JavaScript
			"xhr",           // JavaScript
			"nsurlsession",  // Objective-C & Swift
			"cohttp",        // OCaml
			"guzzle",        // PHP
			"webrequest",    // Powershell
			"restmethod",    // Powershell
			"http.client",   // Python
			"requests",      // Python
			"python3",       // Python
			"HTTPX (Async)", // Python
			"httr",          // R
			"request",       // Unknown
			"http1",         // Unknown
			"http2",         // Unknown
			// "fetch",         // JavaScript & Node.js
			// "axios",         // JavaScript & Node.js
			// "ofetch",        // JavaScript & Node.js
			// "undici",        // Node.js
			// "reqwest",       // Rust
			// "curl",          // Shell & PHP
			// "wget",          // Shell
			// "httpie",        // Shell
		),
		scalar.WithLayout(scalar.LayoutModern),
		scalar.WithHideDarkModeToggle(),
		scalar.WithHideModels(),
		scalar.WithDarkMode(),
		scalar.WithTheme(scalar.ThemePurple),
		scalar.WithOverrideCSS(`
            aside div.flex.items-center:has(a[href*="scalar.com"]) { display: none !important; }
            .scalar-app { font-family: -apple-system, BlinkMacSystemFont, Aptos, "Segoe UI", Roboto, sans-serif; }
        `),
	)

	if err != nil {
		h.Logger.Error("failed to generate Scalar API docs", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate API docs")
	}

	return c.HTML(http.StatusOK, htmlContent)
}

// getFileSystem always uses embed.FS for static assets.
func getFileSystem(embedFS embed.FS, subdir string) http.FileSystem {
	fsys, err := fs.Sub(embedFS, subdir)
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}
