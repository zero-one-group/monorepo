package handler

import (
	"log/slog"
	"net/http"

	"go-modular/modules/user/models"
	"go-modular/modules/user/services"
	"go-modular/pkg/apputils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// HandlerInterface defines the contract for user handlers.
type HandlerInterface interface {
	CreateUser(c echo.Context) error
	ListUsers(c echo.Context) error
	GetUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
}

// Ensure Handler implements HandlerInterface
var _ HandlerInterface = (*Handler)(nil)

// Handler holds dependencies for user handlers.
type Handler struct {
	logger      *slog.Logger
	userService services.UserServiceInterface
	validator   *validator.Validate
}

type HandlerOpts struct {
	Logger      *slog.Logger
	UserService services.UserServiceInterface
}

// NewHandler creates a new Handler instance.
func NewHandler(opts *HandlerOpts) *Handler {
	return &Handler{
		logger:      opts.Logger,
		userService: opts.UserService,
		validator:   validator.New(),
	}
}

// @Summary      Create a new user
// @Description  Creates a new user in the system
// @Tags         User Management
// @Security     BearerAuth
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserCreateRequest  true  "User payload"
// @Success      201   {object}  models.User
// @Failure      400   {object}  map[string]string
// @Router       /api/v1/users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var req models.UserCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}

	user := &models.User{
		DisplayName:     req.Name,
		Email:           req.Email,
		EmailVerifiedAt: nil,
	}

	if err := h.userService.CreateUser(c.Request().Context(), user); err != nil {
		h.logger.Error("Failed to create user", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user, please try again later"})
	}

	return c.JSON(http.StatusCreated, user)
}

// @Summary      List users
// @Description  Retrieves a list of users
// @Tags         User Management
// @Security     BearerAuth
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Produce      json
// @Success      200  {array}   models.User
// @Router       /api/v1/users [get]
func (h *Handler) ListUsers(c echo.Context) error {
	var filter models.FilterUser
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid filter parameters"})
	}

	users, err := h.userService.ListUsers(c.Request().Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to list users", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve users"})
	}

	return c.JSON(http.StatusOK, users)
}

// @Summary      Get user details
// @Description  Retrieves a user by their ID
// @Tags         User Management
// @Security     BearerAuth
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/users/:userId [get]
func (h *Handler) GetUser(c echo.Context) error {
	idStr := c.Param("userId")
	id, err := models.ParseUserID(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID in path must be a valid UUID"})
	}

	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("User not found", slog.String("error", err.Error()))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary      Update user
// @Description  Updates an existing user by ID
// @Tags         User Management
// @Security     BearerAuth
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "User ID"
// @Param        user  body      models.UserCreateRequest  true  "User payload"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /api/v1/users/:userId [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	idStr := c.Param("userId")
	id, err := models.ParseUserID(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID in path must be a valid UUID"})
	}

	var req models.UserCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Validation failed",
			"details": apputils.ValidationErrorsToMap(err, req),
		})
	}

	// Map request fields to user model
	user := &models.User{
		ID:              id,
		DisplayName:     req.Name,
		Email:           req.Email,
		EmailVerifiedAt: nil,
	}

	if err := h.userService.UpdateUser(c.Request().Context(), user); err != nil {
		h.logger.Error("Failed to update user", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user, please try again later"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User updated successfully"})
}

// @Summary      Delete user
// @Description  Deletes a user by ID
// @Tags         User Management
// @Security     BearerAuth
// @Param        Authorization  header    string                      true  "Bearer {token}"
// @Param        id   path  string  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/users/:userId [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	idStr := c.Param("userId")
	id, err := models.ParseUserID(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID in path must be a valid UUID"})
	}

	if err := h.userService.DeleteUser(c.Request().Context(), id); err != nil {
		h.logger.Error("Failed to delete user", slog.String("error", err.Error()))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
