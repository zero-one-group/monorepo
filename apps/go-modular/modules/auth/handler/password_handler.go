package handler

import (
	"log/slog"
	"net/http"

	"go-modular/modules/auth/models"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

// @Summary      Set user password
// @Description  Sets a new password for a user
// @Tags         Auth - User Password
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
// @Tags         Auth - User Password
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
