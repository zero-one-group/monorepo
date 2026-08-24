package server

import "github.com/labstack/echo/v4"

// APIV1 is the /api/v1 route group shared by feature Modules.
type APIV1 struct {
	Root *echo.Group
}
