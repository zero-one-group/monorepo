# Observability

Distributed tracing for this service is built on [OpenTelemetry compile-time
instrumentation](https://opentelemetry.io/blog/2026/go-compile-time-instrumentation-v1/)
(`otelc`), with two library-specific tracers filling the gaps that compile-time
instrumentation cannot cover on its own.

## Quick start

```sh
# 1. Start Jaeger (UI on http://localhost:16686)
pnpm run compose:instrumented

# 2. Install the instrumentation toolchain (once)
moon run go-modular:install-otelc

# 3. Enable telemetry in .env
#    OTEL_ENABLE_TELEMETRY=true

# 4. Build and run the instrumented binary
moon run go-modular:start-otel
```

Send a request, then open <http://localhost:16686>, pick the `go-modular`
service and hit **Find Traces**.

## How the three layers fit together

Compile-time instrumentation alone does **not** produce a complete trace for
this service. It ships an instrumentation for `net/http`, but not for Echo or
pgx, so out of the box every request collapses into a single span named `GET`
with no database activity attached. Three layers are therefore combined:

| Layer | Provides | Why it is needed |
| --- | --- | --- |
| `otelc` (compile-time) | SDK bootstrap, `net/http` client+server spans, `log/slog` correlation, Go runtime metrics | Installs and configures the SDK with no code changes |
| `otelecho` middleware | Renames the span to include the matched route | `otelc` hooks `net/http` and cannot see Echo's routing table |
| `otelpgx` tracer | `pool.acquire`, `prepare` and `query` spans | Only `database/sql` is auto-instrumented; this service uses pgx |

A verified trace for `POST /api/v1/auth/signin/email` looks like this:

```
POST                                    (otelc, net/http)
└── POST /api/v1/auth/signin/email      (otelecho, http.route)
    ├── pool.acquire                    (otelpgx)
    ├── prepare                         (otelpgx, db.query.text)
    └── query                           (otelpgx, db.query.text)
```

The outer bare `POST` span is the `net/http` server span; `otelecho` nests the
route-named span beneath it rather than renaming it in place.

## Building

`otelc` wraps the Go toolchain — everything after `go` is forwarded unchanged:

```sh
otelc go build -o build/release/go-modular ./cmd/
```

Two properties are worth knowing:

- **No source or manifest changes.** `otelc` resolves the OpenTelemetry SDK in
  an isolated build sandbox, so `go.mod` and `go.sum` are left untouched. It may
  report `Bumped dependency ...` lines; those apply to the sandbox only.
- **A plain `go build` still works.** It produces a binary with no SDK, where
  the global tracer provider is the no-op implementation. The `OTEL_*` variables
  are then ignored and the tracing code paths cost nothing.

The `build` / `start` tasks remain uninstrumented. Use `build-otel` /
`start-otel` for tracing.

## Configuration

Configured entirely through standard `OTEL_*` environment variables (see
`.env.example`).

| Variable | Purpose |
| --- | --- |
| `OTEL_ENABLE_TELEMETRY` | Application switch for the Echo and pgx tracers. Does **not** install the SDK. |
| `OTEL_SERVICE_NAME` | Service name shown in Jaeger. |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Collector endpoint. |
| `OTEL_EXPORTER_OTLP_PROTOCOL` | `http/protobuf` (port 4318) or `grpc` (port 4317). |
| `OTEL_TRACES_SAMPLER` / `_ARG` | Sampling. `always_on` is convenient locally. |
| `OTEL_METRICS_EXPORTER` / `OTEL_LOGS_EXPORTER` | Set to `none`; Jaeger ingests traces only. |

### Endpoint and protocol must agree

Jaeger exposes OTLP on **two** ports and the protocol must match the port:

- `4317` — gRPC (`OTEL_EXPORTER_OTLP_PROTOCOL=grpc`)
- `4318` — HTTP (`OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf`)

A mismatch fails silently: the process runs normally and no spans arrive. If
Jaeger shows no service, check this pairing first.

## Adding custom spans

Auto-instrumentation covers transport and database boundaries. For
business-level detail use `internal/observer/tracer`, which reads the globally
registered provider and is a no-op in an uninstrumented build:

```go
func (s *Service) Register(ctx context.Context, email string) error {
    ctx, span := tracer.Start(ctx, "user.Register")
    defer span.End()

    tracer.AddAttributes(ctx, attribute.String("user.email", email))

    if err := s.repo.Insert(ctx, email); err != nil {
        tracer.RecordError(ctx, err) // marks the span failed
        return err
    }
    return nil
}
```

Pass `ctx` down the call chain — that is what makes child spans nest correctly.
`tracer.TraceID(ctx)` returns the current trace id for correlating a log line or
an error response with Jaeger.

## Troubleshooting

**No service listed in Jaeger.** The binary was probably built with a plain
`go build`; confirm with `moon run go-modular:build-otel`. Otherwise check the
endpoint/protocol pairing above.

**Traces contain only a bare `GET` span.** `OTEL_ENABLE_TELEMETRY` is not
`true`, so the Echo and pgx tracers are inactive.

**No database spans.** Same cause — the otelpgx tracer is installed only when
telemetry is enabled at the time the pool is created.
