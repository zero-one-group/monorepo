package handler

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"go-modular/modules/auth/models"
	"go-modular/pkg/apputils"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
)

// @Summary      Create refresh token
// @Description  Creates a new refresh token
// @Tags         Auth - Session Management
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateRefreshTokenRequest  true  "Refresh token payload"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /api/v1/auth/refresh-token [post]
func (h *Handler) CreateRefreshToken(c echo.Context) error {
	var req models.CreateRefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}

	userID, err := uuid.FromString(req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
	}

	var sessionIDPtr *uuid.UUID
	if req.SessionID != nil {
		sid, err := uuid.FromString(*req.SessionID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid session_id"})
		}
		sessionIDPtr = &sid
	}

	var ipPtr *net.IP
	if req.IPAddress != nil {
		ip := net.ParseIP(*req.IPAddress)
		if ip == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ip_address"})
		}
		ipPtr = &ip
	}

	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expires_at"})
	}

	refreshToken := &models.RefreshToken{
		UserID:    userID,
		SessionID: sessionIDPtr,
		TokenHash: []byte(req.TokenHash),
		IPAddress: ipPtr,
		UserAgent: req.UserAgent,
		ExpiresAt: expiresAt,
	}

	if err := h.authService.CreateRefreshToken(c.Request().Context(), refreshToken); err != nil {
		h.logger.Error("Failed to create refresh token", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create refresh token"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Refresh token created successfully"})
}

// @Summary      Update refresh token
// @Description  Updates an existing refresh token
// @Tags         Auth - Session Management
// @Accept       json
// @Produce      json
// @Param        body  body      models.UpdateRefreshTokenRequest  true  "Refresh token payload"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /api/v1/auth/refresh-token [put]
func (h *Handler) UpdateRefreshToken(c echo.Context) error {
	var req models.UpdateRefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}

	tokenID, err := uuid.FromString(req.TokenID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid token_id"})
	}

	var ipPtr *net.IP
	if req.IPAddress != nil {
		ip := net.ParseIP(*req.IPAddress)
		if ip == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ip_address"})
		}
		ipPtr = &ip
	}

	var revokedAt *time.Time
	if req.RevokedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.RevokedAt)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid revoked_at"})
		}
		revokedAt = &t
	}

	var revokedBy *uuid.UUID
	if req.RevokedBy != nil {
		id, err := uuid.FromString(*req.RevokedBy)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid revoked_by"})
		}
		revokedBy = &id
	}

	refreshToken := &models.RefreshToken{
		ID:        tokenID,
		IPAddress: ipPtr,
		UserAgent: req.UserAgent,
		RevokedAt: revokedAt,
		RevokedBy: revokedBy,
	}

	if err := h.authService.UpdateRefreshToken(c.Request().Context(), refreshToken); err != nil {
		h.logger.Error("Failed to update refresh token", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update refresh token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Refresh token updated successfully"})
}

// @Summary      Get refresh token
// @Description  Retrieves a refresh token by its ID
// @Tags         Auth - Session Management
// @Produce      json
// @Param        tokenId  path      string  true  "Token ID"
// @Success      200      {object}  models.RefreshToken
// @Failure      400      {object}  map[string]string
// @Router       /api/v1/auth/refresh-token/{tokenId} [get]
func (h *Handler) GetRefreshToken(c echo.Context) error {
	tokenIDStr := c.Param("tokenId")
	tokenID, err := uuid.FromString(tokenIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid token_id"})
	}

	token, err := h.authService.GetRefreshToken(c.Request().Context(), tokenID)
	if err != nil {
		h.logger.Error("Failed to get refresh token", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get refresh token"})
	}
	if token == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Refresh token not found"})
	}

	return c.JSON(http.StatusOK, token)
}

// @Summary      Delete refresh token
// @Description  Deletes a refresh token by its ID
// @Tags         Auth - Session Management
// @Produce      json
// @Param        tokenId  path      string  true  "Token ID"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Router       /api/v1/auth/refresh-token/{tokenId} [delete]
func (h *Handler) DeleteRefreshToken(c echo.Context) error {
	tokenIDStr := c.Param("tokenId")
	tokenID, err := uuid.FromString(tokenIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid token_id"})
	}

	if err := h.authService.DeleteRefreshToken(c.Request().Context(), tokenID); err != nil {
		h.logger.Error("Failed to delete refresh token", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete refresh token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Refresh token deleted successfully"})
}
