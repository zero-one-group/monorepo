package handler

import (
	"net/http"

	"go-modular/modules/auth/models"
	"go-modular/pkg/apputils"

	"github.com/labstack/echo/v4"
	"go-modular/modules/user/repository"
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}
	err := h.authService.InitiateEmailVerification(ctx, req.Email)
	if err != nil {
		if err.Error() == "user not found" || err == repository.ErrNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error":   "User with this email is not registered",
				"details": err.Error(),
			})
		}
		if err.Error() == "token still valid" {
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"error": "A verification token is still valid. Please check your email or try again later.",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to initiate email verification",
			"details": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Verification email sent if the email is registered",
	})
}

// @Summary      Validate email verification
// @Description  Validates the email verification token for the user
// @Tags         Auth - Verification
// @Accept       json
// @Produce      json
// @Param        body  body      models.ValidateEmailVerificationRequest  true  "Validate verification payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/auth/verification/email/validate [post]
func (h *Handler) ValidateEmailVerification(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.ValidateEmailVerificationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}
	ok, err := h.authService.ValidateEmailVerification(ctx, req.Email, req.Token)
	if err != nil {
		if err.Error() == "user not found" || err == repository.ErrNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error":   "User with this email is not registered",
				"details": err.Error(),
			})
		}
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error":   "Invalid or expired token",
			"details": err.Error(),
		})
	}
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Invalid or expired token",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Email successfully verified",
	})
}

// @Summary      Revoke email verification
// @Description  Revokes (deletes) the email verification token for the user
// @Tags         Auth - Verification
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}
	if err := h.authService.RevokeEmailVerification(ctx, req.Token); err != nil {
		if err == repository.ErrNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error":   "Token not found or already revoked",
				"details": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to revoke verification token",
			"details": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Verification token revoked",
	})
}

// @Summary      Resend email verification
// @Description  Generates and sends a new verification token to the user's email, or resends if still valid
// @Tags         Auth - Verification
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
	var req models.ResendEmailVerificationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}
	err := h.authService.ResendEmailVerification(ctx, req.Email)
	if err != nil {
		if err.Error() == "user not found" || err == repository.ErrNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error":   "User with this email is not registered",
				"details": err.Error(),
			})
		}
		if err.Error() == "token still valid" {
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"error": "A verification token is still valid. Please check your email or try again later.",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to resend verification email",
			"details": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Verification email resent if the email is registered",
	})
}
