package handler

import (
	"log/slog"
	"net/http"

	"go-modular/modules/auth/models"
	"go-modular/modules/auth/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

// HandlerInterface defines the contract for auth handlers.
type HandlerInterface interface {
	SetUserPassword(c echo.Context) error
	UpdateUserPassword(c echo.Context) error
}

// Ensure Handler implements HandlerInterface
var _ HandlerInterface = (*Handler)(nil)

// Handler holds dependencies for auth handlers.
type Handler struct {
	logger      *slog.Logger
	authService services.AuthServiceInterface
	validator   *validator.Validate
}

type HandlerOpts struct {
	Logger      *slog.Logger
	AuthService services.AuthServiceInterface
}

// NewHandler creates a new Handler instance.
func NewHandler(opts *HandlerOpts) *Handler {
	return &Handler{
		logger:      opts.Logger,
		authService: opts.AuthService,
		validator:   validator.New(),
	}
}

// @Summary      Set user password
// @Description  Sets a new password for a user
// @Tags         Auth - Password
// @Accept       json
// @Produce      json
// @Param        body  body      models.SetPasswordRequest  true  "Password payload"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /api/v1/auth/password [post]
func (h *Handler) SetUserPassword(c echo.Context) error {
	var req models.SetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": validationErrorsToMap(err),
		})
	}

	userID, err := uuid.FromString(req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
	}

	userPassword := &models.UserPassword{
		UserID:       userID,
		PasswordHash: req.Password,
	}

	if err := h.authService.SetUserPassword(c.Request().Context(), userPassword); err != nil {
		h.logger.Error("Failed to set user password", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to set password"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Password set successfully"})
}

// @Summary      Update user password
// @Description  Updates an existing user's password
// @Tags         Auth - Password
// @Accept       json
// @Produce      json
// @Param        userId  path      string                      true  "User ID"
// @Param        body    body      models.UpdatePasswordRequest true  "Password payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Router       /api/v1/auth/password/{userId} [put]
func (h *Handler) UpdateUserPassword(c echo.Context) error {
	userIDStr := c.Param("userId")
	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
	}

	var req models.UpdatePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": validationErrorsToMap(err),
		})
	}

	// Validasi current password sebelum update
	if err := h.authService.UpdateUserPassword(
		c.Request().Context(),
		userID,
		req.CurrentPassword,
		req.NewPassword,
	); err != nil {
		h.logger.Error("Failed to update user password", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password updated successfully"})
}

// Helper to convert validator.ValidationErrors to a readable map
func validationErrorsToMap(err error) map[string]string {
	errs := map[string]string{}
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()
			var msg string
			switch tag {
			case "required":
				msg = "This field is required"
			case "uuid":
				msg = "Must be a valid UUID"
			case "min":
				msg = "Minimum length is " + fe.Param()
			case "eqfield":
				msg = "Must match " + fe.Param()
			default:
				msg = "Invalid value"
			}
			errs[field] = msg
		}
	} else {
		errs["error"] = err.Error()
	}
	return errs
}
