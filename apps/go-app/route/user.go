package router

import (
	"go-app/domain"
	"go-app/service"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

func RegisterUserRoutes(e *echo.Group, svc *service.UserService) {
	e.POST("/users", createUserHandler(svc))
	e.GET("/users/:id", getUserHandler(svc))
	e.PUT("/users/:id", updateUserHandler(svc))
	e.DELETE("/users/:id", deleteUserHandler(svc))
}

func createUserHandler(svc *service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req domain.User
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest,
				map[string]string{"error": "invalid request body"})
		}
		user, err := svc.CreateUser(c.Request().Context(), &req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusCreated, user)
	}
}

func getUserHandler(svc *service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(
			c.Request().Context(),
			"RouteUser.GetUser",
		)
		defer span.Finish()

		rawID := c.Param("id")
		span.SetBaggageItem("user_id", rawID)

		id, err := uuid.Parse(rawID)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{"error": "invalid user id"},
			)
		}
		span.SetTag("user.id", id.String())

		user, err := svc.GetUser(ctx, id)
		if err != nil {
			if err == domain.ErrUserNotFound {
				return c.JSON(
					http.StatusNotFound,
					map[string]string{"error": "user not found"},
				)
			}
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"error": err.Error()},
			)
		}

		return c.JSON(http.StatusOK, user)
	}
}

func updateUserHandler(svc *service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest,
				map[string]string{"error": "invalid user id"})
		}
		var req domain.User
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest,
				map[string]string{"error": "invalid request body"})
		}
		user, err := svc.UpdateUser(c.Request().Context(), id, &req)
		if err != nil {
			if err == domain.ErrUserNotFound {
				return c.JSON(http.StatusNotFound,
					map[string]string{"error": "user not found"})
			}
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, user)
	}
}

func deleteUserHandler(svc *service.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest,
				map[string]string{"error": "invalid user id"})
		}
		if err := svc.DeleteUser(c.Request().Context(), id); err != nil {
			if err == domain.ErrUserNotFound {
				return c.JSON(http.StatusNotFound,
					map[string]string{"error": "user not found"})
			}
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"error": err.Error()})
		}
		return c.NoContent(http.StatusNoContent)
	}
}
