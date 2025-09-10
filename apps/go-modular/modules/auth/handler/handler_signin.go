package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-modular/modules/auth/models"
)

// SignInWithEmail godoc
// @Summary      Sign in with email
// @Description  Authenticates user using email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.SignInWithEmailRequest  true  "Sign in payload"
// @Success      200   {object}  models.SignInResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /api/v1/auth/signin/email [post]
func (h *Handler) SignInWithEmail(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.SignInWithEmailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": validationErrorsToMap(err),
		})
	}

	authedUser, err := h.authService.SignInWithEmail(ctx, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	return c.JSON(http.StatusOK, authedUser)
}

// SignInWithUsername godoc
// @Summary      Sign in with username
// @Description  Authenticates user using username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.SignInWithUsernameRequest  true  "Sign in payload"
// @Success      200   {object}  models.SignInResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /api/v1/auth/signin/username [post]
func (h *Handler) SignInWithUsername(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.SignInWithUsernameRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": validationErrorsToMap(err),
		})
	}

	authedUser, err := h.authService.SignInWithUsername(ctx, req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
	}

	return c.JSON(http.StatusOK, authedUser)
}
