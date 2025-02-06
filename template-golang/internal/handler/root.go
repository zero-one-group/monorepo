package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"{{ package_id | kebab_case }}/internal/model"
)

func RootHandler(c echo.Context) error {
	response := model.SuccessResponse[map[string]string]{
		BaseAPIResponse: model.BaseAPIResponse{
			Status:  http.StatusOK,
			Success: true,
			Message: "Nothing to see here",
		},
		Data: map[string]string{},
	}
	return c.JSON(http.StatusOK, response)
}

// TODO: use health check package like: https://github.com/alexliesenfeld/health
func HealthCheckHandler(c echo.Context) error {
	response := model.HealthCheckResponse{
		BaseAPIResponse: model.BaseAPIResponse{
			Status:  http.StatusOK,
			Success: true,
			Message: "All is well",
		},
		Data: model.HealthCheckData{
			Uptime:    0000000000,
			Timestamp: time.Now().Unix(),
		},
	}
	return c.JSON(http.StatusOK, response)
}
