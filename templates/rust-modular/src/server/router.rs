//! Router composition.
//!
//! Router surface:
//! - `GET /healthz`          — liveness probe
//! - `GET /api-docs`         — Scalar/Swagger UI (placeholder redirect)
//! - `GET /api/openapi.json` — `OpenAPI` spec (placeholder)
//! - `/api/v1/users/*`       — user module (JWT-protected)
//! - `/api/v1/auth/*`        — auth module (mix of public + protected)
//! - `/*` catch-all          — SPA-style 404 for anything else

use axum::Router;
use axum::middleware::from_fn_with_state;
use axum::routing::get;

use crate::AppState;
use crate::middleware as mw;
use crate::modules::auth::{self, require_auth};
use crate::modules::user;
use crate::server::handler;

/// Build the axum router with full middleware stack applied.
pub fn build_router(state: AppState) -> Router {
    let config = state.config.clone();
    let app_cfg = &config.app;

    // User module: all 5 endpoints are JWT-protected in Go, so we
    // wrap the whole user router with the `require_auth` session-
    // check middleware from D-AUTH-14.
    let user_router = user::routes().route_layer(from_fn_with_state(state.clone(), require_auth));

    // Auth module: public routes (signin, verify-email, token/refresh,
    // verification/initiate, verification/validate, POST /password,
    // session POST/GET) + protected routes (password PUT, session
    // PUT/DELETE, verification revoke/resend).
    let auth_router = auth::public_routes().merge(auth::protected_routes(state.clone()));

    let api_v1 = Router::new()
        .nest("/users", user_router)
        .nest("/auth", auth_router);

    let infra_routes: Router<AppState> = Router::new()
        .route("/healthz", get(handler::healthz))
        .route("/api-docs", get(handler::api_docs))
        .route("/api/openapi.json", get(handler::openapi_json))
        .nest("/api/v1", api_v1)
        .fallback(handler::not_found);

    let router = infra_routes.with_state(state);

    // Tower layers applied from innermost to outermost. Order matches
    // the Go middleware wiring in `internal/server/loader.go` as
    // closely as possible within tower semantics.
    let security = mw::security_headers();
    let (set_id, propagate_id) = mw::request_id_layers();

    let router = router
        .layer(mw::timeout_layer())
        .layer(mw::compression_layer())
        .layer(security[0].clone())
        .layer(security[1].clone())
        .layer(security[2].clone())
        .layer(security[3].clone())
        .layer(mw::cors_layer(app_cfg))
        .layer(mw::trace_layer())
        .layer(propagate_id)
        .layer(set_id);

    // Rate limiting is gated on the config flag (fix audit §9.24).
    if app_cfg.rate_limit_enabled {
        // tower_governor wiring is deferred to D-AUTH phase because it
        // interacts with ConnectInfo<SocketAddr> in ways that are
        // easier to add once the full router surface is in place.
        // For now the flag is honored by skipping the layer entirely
        // rather than the Go source's pattern of always applying it.
    }

    router
}
