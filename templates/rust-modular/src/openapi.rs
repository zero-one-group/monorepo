//! `OpenAPI` surface (D-DOC-1).
//!
//! `ApiDoc::openapi()` returns the complete generated spec,
//! enumerating all 19 API endpoints + referenced schemas.
//! `server::handler::openapi_json` serializes this and returns it
//! as `application/json`.

// utoipa's `#[derive(OpenApi)]` macro expansion trips
// `clippy::needless_for_each` inside its generated code. The
// warning is on macro-authored code, not ours, so suppress it
// at the module level.
#![allow(clippy::needless_for_each)]

use utoipa::OpenApi;
use utoipa::openapi::security::{HttpAuthScheme, HttpBuilder, SecurityScheme};

use crate::domain::MessageResponse;
use crate::domain::response::ErrorBody;
use crate::modules::auth::handler as auth_handler;
use crate::modules::auth::models::{
    AuthenticatedUser, OneTimeToken, RefreshToken, Session, UserPassword,
};
use crate::modules::auth::schema::{
    CreateSessionRequest, InitiateEmailVerificationRequest, NeutralResponse,
    ResendEmailVerificationRequest, RevokeEmailVerificationRequest, SetPasswordRequest,
    SignInWithEmailRequest, SignInWithUsernameRequest, TokenRefreshRequest, UpdatePasswordRequest,
    UpdateSessionRequest, ValidateEmailVerificationRequest,
};
use crate::modules::user::handler as user_handler;
use crate::modules::user::models::{User, UserCreateRequest, UserMetadata};

#[derive(OpenApi)]
#[openapi(
    info(
        title = "{{ package_name | kebab_case }}",
        version = "0.1.0",
        description = "Rust port of the Zero One Group {{ package_name | kebab_case }} service (Phase D). \
                       14 auth endpoints + 5 user endpoints with 8 corrected-port design fixes."
    ),
    paths(
        // User module (5)
        user_handler::create_user,
        user_handler::list_users,
        user_handler::get_user,
        user_handler::update_user,
        user_handler::delete_user,
        // Auth public (9)
        auth_handler::sign_in_with_email,
        auth_handler::sign_in_with_username,
        auth_handler::verify_email_by_link,
        auth_handler::rotate_refresh_token,
        auth_handler::initiate_email_verification,
        auth_handler::validate_email_verification,
        auth_handler::set_password,
        auth_handler::create_session,
        auth_handler::get_session,
        // Auth protected (5)
        auth_handler::update_password,
        auth_handler::update_session,
        auth_handler::delete_session,
        auth_handler::revoke_email_verification,
        auth_handler::resend_email_verification,
    ),
    components(schemas(
        // Response shapes
        User,
        UserMetadata,
        MessageResponse,
        ErrorBody,
        AuthenticatedUser,
        Session,
        RefreshToken,
        OneTimeToken,
        UserPassword,
        NeutralResponse,
        // Request bodies
        UserCreateRequest,
        SetPasswordRequest,
        UpdatePasswordRequest,
        CreateSessionRequest,
        UpdateSessionRequest,
        SignInWithEmailRequest,
        SignInWithUsernameRequest,
        TokenRefreshRequest,
        InitiateEmailVerificationRequest,
        ValidateEmailVerificationRequest,
        RevokeEmailVerificationRequest,
        ResendEmailVerificationRequest,
    )),
    tags(
        (name = "User Management", description = "User CRUD endpoints"),
        (name = "Authentication", description = "Signin + token rotation"),
        (name = "Password", description = "Set + update password (ownership-checked)"),
        (name = "Sessions", description = "Session CRUD + revocation"),
        (name = "Email Verification", description = "Initiate + validate + resend"),
    ),
    modifiers(&SecurityAddon),
)]
pub struct ApiDoc;

/// Registers the `bearer_auth` security scheme referenced by
/// individual endpoints via `security(("bearer_auth" = []))`.
struct SecurityAddon;

impl utoipa::Modify for SecurityAddon {
    fn modify(&self, openapi: &mut utoipa::openapi::OpenApi) {
        let components = openapi
            .components
            .as_mut()
            .expect("utoipa-generated components");
        components.add_security_scheme(
            "bearer_auth",
            SecurityScheme::Http(
                HttpBuilder::new()
                    .scheme(HttpAuthScheme::Bearer)
                    .bearer_format("JWT")
                    .build(),
            ),
        );
    }
}

#[cfg(test)]
mod tests {
    use super::ApiDoc;
    use utoipa::OpenApi;

    /// Assert the spec is valid JSON and contains all 19 paths.
    #[test]
    fn openapi_spec_has_all_19_paths() {
        let doc = ApiDoc::openapi();
        let value: serde_json::Value = serde_json::to_value(&doc).expect("spec serializes to JSON");

        let paths = value
            .get("paths")
            .and_then(|p| p.as_object())
            .expect("paths object present");

        // Note: utoipa groups path items by URL, and each URL can
        // carry multiple HTTP methods. Our 19 endpoints span 16
        // unique URLs (POST /session + GET/PUT /session collapse
        // onto a single path entry; same for /users/{userId}).
        //
        // Assert each expected URL is present, then count the total
        // number of operations across all URLs.
        let expected_urls = [
            "/api/v1/users",
            "/api/v1/users/{userId}",
            "/api/v1/auth/signin/email",
            "/api/v1/auth/signin/username",
            "/api/v1/auth/verify-email",
            "/api/v1/auth/token/refresh",
            "/api/v1/auth/verification/email/initiate",
            "/api/v1/auth/verification/email/validate",
            "/api/v1/auth/verification/email/revoke",
            "/api/v1/auth/verification/email/resend",
            "/api/v1/auth/password",
            "/api/v1/auth/password/{userId}",
            "/api/v1/auth/session",
            "/api/v1/auth/session/{sessionId}",
        ];
        for url in expected_urls {
            assert!(
                paths.contains_key(url),
                "missing path in spec: {url}\nfull paths: {:?}",
                paths.keys().collect::<Vec<_>>()
            );
        }

        // Count operations (method + url pairs) across all paths.
        let operations: usize = paths
            .values()
            .map(|path_item| {
                path_item.as_object().map_or(0, |ops| {
                    ops.keys()
                        .filter(|k| {
                            matches!(
                                k.as_str(),
                                "get" | "post" | "put" | "delete" | "patch" | "head"
                            )
                        })
                        .count()
                })
            })
            .sum();
        assert_eq!(
            operations, 19,
            "expected 19 total operations across all paths, got {operations}"
        );
    }

    #[test]
    fn openapi_spec_has_info_and_security_scheme() {
        let doc = ApiDoc::openapi();
        let value: serde_json::Value = serde_json::to_value(&doc).expect("spec serializes");

        assert_eq!(value["info"]["title"], "{{ package_name | kebab_case }}");
        assert_eq!(
            value["components"]["securitySchemes"]["bearer_auth"]["scheme"],
            "bearer"
        );
    }
}
