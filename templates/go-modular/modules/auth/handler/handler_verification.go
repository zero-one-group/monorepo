package handler

import (
	"net/http"
	"net/url"

	"{{ package_name | kebab_case }}/modules/auth/models"
	"{{ package_name | kebab_case }}/pkg/apputils"

	"github.com/labstack/echo/v4"
	"{{ package_name | kebab_case }}/modules/user/repository"
)

// @Summary      Initiate email verification
// @Description  Generates and sends a verification token to the user's email
// @Tags         Auth - Verification
// @Accept       json
// @Produce      json
// @Param        body  body      models.InitiateEmailVerificationRequest  true  "Initiate verification payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/auth/verification/email/initiate [post]
func (h *Handler) InitiateEmailVerification(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.InitiateEmailVerificationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}

	// pass optional redirect_to to service so it will be stored in token metadata
	err := h.authService.InitiateEmailVerification(ctx, req.Email, req.RedirectTo)
	if err != nil {
		if err.Error() == "user not found" || err == repository.ErrNotFound {
			return c.JSON(http.StatusNotFound, map[string]any{
				"error":   "User with this email is not registered",
				"details": err.Error(),
			})
		}
		if err.Error() == "token still valid" {
			return c.JSON(http.StatusConflict, map[string]any{
				"error": "A verification token is still valid. Please check your email or try again later.",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   "Failed to initiate email verification",
			"details": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"message": "Verification email sent if the email is registered",
	})
}

// @Summary      Validate email verification
// @Description  Validates the email verification token for the user (token-only payload)
// @Tags         Auth - Verification
// @Accept       json
// @Produce      json
// @Param        body  body      map[string]string  true  "Validate verification payload (json: {\"token\":\"...\"})"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/auth/verification/email/validate [post]
func (h *Handler) ValidateEmailVerification(c echo.Context) error {
	ctx := c.Request().Context()

	// Accept token-only payload to match verification links that include only the token.
	var req struct {
		Token string `json:"token" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}

	ok, err := h.authService.ValidateEmailVerification(ctx, req.Token)
	if err != nil {
		// treat any service error as unauthorized/invalid token for this endpoint
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"error":   "Invalid or expired token",
			"details": err.Error(),
		})
	}
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"error": "Invalid or expired token",
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"message": "Email successfully verified",
	})
}

// @Summary      Revoke email verification
// @Description  Revokes (deletes) the email verification token for the user
// @Tags         Auth - Verification
// @Security     BearerAuth
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Accept       json
// @Produce      json
// @Param        body  body      models.RevokeEmailVerificationRequest  true  "Revoke verification payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/auth/verification/email/revoke [post]
func (h *Handler) RevokeEmailVerification(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.RevokeEmailVerificationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}
	if err := h.authService.RevokeEmailVerification(ctx, req.Token); err != nil {
		if err == repository.ErrNotFound {
			return c.JSON(http.StatusNotFound, map[string]any{
				"error":   "Token not found or already revoked",
				"details": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   "Failed to revoke verification token",
			"details": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"message": "Verification token revoked",
	})
}

// @Summary      Resend email verification
// @Description  Generates and sends a new verification token to the user's email, or resends if still valid
// @Tags         Auth - Verification
// @Security     BearerAuth
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Accept       json
// @Produce      json
// @Param        body  body      models.ResendEmailVerificationRequest  true  "Resend verification payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/auth/verification/email/resend [post]
func (h *Handler) ResendEmailVerification(c echo.Context) error {
	ctx := c.Request().Context()

	// allow optional redirect_to here (store in metadata); bind into local struct to avoid requiring model change
	var req struct {
		Email      string `json:"email" validate:"required,email"`
		RedirectTo string `json:"redirect_to,omitempty" validate:"omitempty,url"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}
	err := h.authService.ResendEmailVerification(ctx, req.Email, req.RedirectTo)
	if err != nil {
		if err.Error() == "user not found" || err == repository.ErrNotFound {
			return c.JSON(http.StatusNotFound, map[string]any{
				"error":   "User with this email is not registered",
				"details": err.Error(),
			})
		}
		if err.Error() == "token still valid" {
			return c.JSON(http.StatusConflict, map[string]any{
				"error": "A verification token is still valid. Please check your email or try again later.",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   "Failed to resend verification email",
			"details": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"message": "Verification email resent if the email is registered",
	})
}

// @Summary      Validate email via link
// @Description  Validates an email verification token provided as query parameter `token`.
//
//	If `redirect_to` is provided and valid, the handler redirects (302 Found)
//	to that URL and appends query params `verified=true|false` and optional `error`.
//	If `redirect_to` is not provided the endpoint returns JSON (or 401 on invalid token).
//
// @Tags         Auth - Verification
// @Param        token       query    string  true   "Verification token"
// @Param        redirect_to query    string  false  "Optional absolute URL to redirect after verification (will receive `verified` and optional `error` query params)" format(url)
// @Success      200         {object} map[string]interface{} "Email successfully verified (JSON)"
// @Success      302         {string} string                 "Redirect to `redirect_to` with verification result"
// @Failure      400         {object} map[string]interface{} "Bad request (missing token or invalid redirect_to)"
// @Failure      401         {object} map[string]interface{} "Invalid or expired token"
// @Failure      500         {object} map[string]interface{} "Server error"
// @Router       /api/v1/auth/verify-email [get]
func (h *Handler) ValidateEmailVerificationByLink(c echo.Context) error {
	ctx := c.Request().Context()
	token := c.QueryParam("token")
	redirectTo := c.QueryParam("redirect_to")

	if token == "" {
		if redirectTo != "" {
			u, err := url.Parse(redirectTo)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid redirect_to URL"})
			}
			q := u.Query()
			q.Set("verified", "false")
			q.Set("error", "token_required")
			u.RawQuery = q.Encode()
			return c.Redirect(http.StatusFound, u.String())
		}
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "token query parameter is required",
		})
	}

	ok, err := h.authService.ValidateEmailVerification(ctx, token)
	if redirectTo != "" {
		u, perr := url.Parse(redirectTo)
		if perr != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid redirect_to URL"})
		}
		q := u.Query()
		if err != nil || !ok {
			q.Set("verified", "false")
			if err != nil {
				q.Set("error", errString(err))
			}
			u.RawQuery = q.Encode()
			return c.Redirect(http.StatusFound, u.String())
		}
		q.Set("verified", "true")
		u.RawQuery = q.Encode()
		return c.Redirect(http.StatusFound, u.String())
	}

	// No redirect requested — return JSON REST responses
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"error":   "Invalid or expired token",
			"details": errString(err),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"message": "Email successfully verified",
	})
}

// helper to safely extract error string (avoid nil panic)
func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
