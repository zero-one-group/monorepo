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

// @Summary      Create session
// @Description  Creates a new session
// @Tags         Auth - User Session
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateSessionRequest  true  "Session payload"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /api/v1/auth/session [post]
func (h *Handler) CreateSession(c echo.Context) error {
	var req models.CreateSessionRequest
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

	session := &models.Session{
		UserID:            userID,
		TokenHash:         req.TokenHash,
		UserAgent:         req.UserAgent,
		DeviceName:        req.DeviceName,
		DeviceFingerprint: req.DeviceFingerprint,
		IPAddress:         ipPtr,
		ExpiresAt:         expiresAt,
	}

	if err := h.authService.CreateSession(c.Request().Context(), session); err != nil {
		h.logger.Error("Failed to create session", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Session created successfully"})
}

// @Summary      Update session
// @Description  Updates an existing session
// @Tags         Auth - User Session
// @Accept       json
// @Produce      json
// @Param        body  body      models.UpdateSessionRequest  true  "Session payload"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /api/v1/auth/session [put]
func (h *Handler) UpdateSession(c echo.Context) error {
	var req models.UpdateSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}

	sessionID, err := uuid.FromString(req.SessionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid session_id"})
	}

	var ipPtr *net.IP
	if req.IPAddress != nil {
		ip := net.ParseIP(*req.IPAddress)
		if ip == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ip_address"})
		}
		ipPtr = &ip
	}

	var refreshedAt, revokedAt *time.Time
	if req.RefreshedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.RefreshedAt)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid refreshed_at"})
		}
		refreshedAt = &t
	}
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

	session := &models.Session{
		ID:                sessionID,
		UserAgent:         req.UserAgent,
		DeviceName:        req.DeviceName,
		DeviceFingerprint: req.DeviceFingerprint,
		IPAddress:         ipPtr,
		RefreshedAt:       refreshedAt,
		RevokedAt:         revokedAt,
		RevokedBy:         revokedBy,
	}

	if err := h.authService.UpdateSession(c.Request().Context(), session); err != nil {
		h.logger.Error("Failed to update session", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update session"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Session updated successfully"})
}

// @Summary      Get session
// @Description  Retrieves a session by its ID
// @Tags         Auth - User Session
// @Produce      json
// @Param        sessionId  path      string  true  "Session ID"
// @Success      200        {object}  models.Session
// @Failure      400        {object}  map[string]string
// @Router       /api/v1/auth/session/:sessionId [get]
func (h *Handler) GetSession(c echo.Context) error {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.FromString(sessionIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid session_id"})
	}

	session, err := h.authService.GetSession(c.Request().Context(), sessionID)
	if err != nil {
		h.logger.Error("Failed to get session", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get session"})
	}
	if session == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Session not found"})
	}

	return c.JSON(http.StatusOK, session)
}

// @Summary      Delete session
// @Description  Deletes a session by its ID
// @Tags         Auth - User Session
// @Produce      json
// @Param        sessionId  path      string  true  "Session ID"
// @Success      200        {object}  map[string]string
// @Failure      400        {object}  map[string]string
// @Router       /api/v1/auth/session/:sessionId [delete]
func (h *Handler) DeleteSession(c echo.Context) error {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.FromString(sessionIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid session_id"})
	}

	if err := h.authService.DeleteSession(c.Request().Context(), sessionID); err != nil {
		h.logger.Error("Failed to delete session", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete session"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Session deleted successfully"})
}
