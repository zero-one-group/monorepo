# Adopt uber/fx for serve and Module composition

Manual wiring in `HTTPServer.Start` / `registerModules` does not scale as a template contract for humans or AI. We adopt `go.uber.org/fx` for the `serve` path only: `internal/app` is the composition root; infrastructure and each feature Module are `fx.Module`s with Lifecycle hooks; feature Modules expose constructors via `fx.go` (no `*Module` facade); Auth provides JWT middleware as a type User route registration consumes; config is loaded by Cobra and `fx.Supply`’d; mailer stays soft-fail. Migrate/seed/generate-config stay outside fx.

## Considered Options

- **Keep manual DI** — simple now, brittle Module graph and lifecycle as the template grows.
- **google/wire** — compile-time graph, no runtime lifecycle; weaker fit for DB/HTTP/mailer start-stop.
- **fx for every Cobra command** — migrate/seed pay for a graph they do not need.

## Consequences

- One PR cutover for User/Auth; explicit Module list in `internal/app` (no auto-discover).
- Fixed Module layout: `handler/`, `services/`, `repository/`, `models/`, root `fx.go`.
