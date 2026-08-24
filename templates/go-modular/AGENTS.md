# Go + Echo Agent Rules

## Project Context

You are an expert Go developer working with Echo.

## Code Style & Structure

### Go Defaults

- Follow Go idioms: short variable names in small scopes, descriptive names for exported identifiers.
- Return errors as the last return value — check immediately with `if err != nil { return err }`. Never discard errors with `_`.
- Use short variable names for short-lived variables (`i`, `n`, `err`) and descriptive names for package-level and long-lived variables.
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)` for traceability.
- Accept interfaces, return structs. Keep interfaces small (1-3 methods).
- Use `context.Context` as the first parameter for functions that do I/O or may be cancelled.
- Prefer returning `(result, error)` tuples over panicking. Reserve `panic` for truly unrecoverable states.
- Organize packages by domain, not by type. Avoid packages named `util`, `common`, or `helpers`.

### Effective Go

- Run gofmt on every save and in CI — the Go community has no style debates; gofmt output is the canonical format.
- Use lowercase, single-word names for packages.
- Provide package documentation in doc.go and godoc comments for all public functions, types, variables, and fields.
- Handle errors explicitly by returning them and wrapping with fmt.Errorf("context: %w", err).
- Pass context.Context as the first parameter in functions that may block or be cancellable.
- Avoid global variables; manage state through function parameters and return values for better testability and clarity.
- Name interfaces using nouns or adjectives (e.g., Reader); use plain names for getters without "Get" prefix (e.g., Config()).
- Group related struct fields logically, such as dependencies followed by configuration.
- Use constructor functions like NewX to initialize structs and validate dependencies.
- Eliminate duplicate code to keep the codebase simple and maintainable.

## Echo

- Initialize the Echo instance with `e := echo.New()`.
- Define routes using methods like `e.GET(path, handler)`, `e.POST(path, handler)`.
- Start the server via `e.Start(addr)` or `e.StartTLS(addr, cert, key)`.
- Access request context in handlers with `c *echo.Context`.
- Organize routes with groups: `g := e.Group("/api")`; then `g.GET("/users", handler)`.
- Handle JSON binding and validation: `var user User; if err := c.Bind(&user); err != nil { ... }`.
- Apply global middleware: `e.Use(middleware.Logger(), middleware.Recover())`.
- Return JSON responses: `c.JSON(http.StatusOK, map[string]interface{}{"data": data})`.
- Use `c.Param("id")` or `c.QueryParam("q")` for path and query parameters.
