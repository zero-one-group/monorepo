package handler

import (
	"errors"
	"net/http"

	"go-modular/modules/auth/models"
	"go-modular/modules/auth/services"
	"go-modular/pkg/apputils"

	"github.com/labstack/echo/v4"
)

// SignInWithEmail godoc
// @Summary      Sign in with email
// @Description  Authenticates user using email and password
// @Tags         Auth - Authentication
// @Accept       json
// @Produce      json
// @Param        body  body      models.SignInWithEmailRequest  true  "Sign in payload"
// @Success      200   {object}  models.SignInResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /api/v1/auth/signin/email [post]
func (h *Handler) SignInWithEmail(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.SignInWithEmailRequest
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

	authedUser, err := h.authService.SignInWithEmail(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Invalid email or password",
				"details": "The email or password you entered is incorrect.",
			})
		case errors.Is(err, services.ErrEmailNotVerified):
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Email is not verified",
				"details": "Please verify your email address before signing in.",
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   "Internal server error",
				"details": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, authedUser)
}

// SignInWithUsername godoc
// @Summary      Sign in with username
// @Description  Authenticates user using username and password
// @Tags         Auth - Authentication
// @Accept       json
// @Produce      json
// @Param        body  body      models.SignInWithUsernameRequest  true  "Sign in payload"
// @Success      200   {object}  models.SignInResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /api/v1/auth/signin/username [post]
func (h *Handler) SignInWithUsername(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.SignInWithUsernameRequest
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

	authedUser, err := h.authService.SignInWithUsername(ctx, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Invalid username or password",
				"details": "The username or password you entered is incorrect.",
			})
		case errors.Is(err, services.ErrEmailNotVerified):
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Email is not verified",
				"details": "Please verify your email address before signing in.",
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   "Internal server error",
				"details": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, authedUser)
}
