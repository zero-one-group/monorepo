package config

import (
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
)

type CloseFunc func() error

func ApplyInstrumentation(e *echo.Echo) CloseFunc {
	var closers []CloseFunc

	closers = append(closers, setupJaeger(e))

	return func() error {
		for i := len(closers) - 1; i >= 0; i-- {
			if err := closers[i](); err != nil {
				return err
			}
		}
		return nil
	}
}

func setupJaeger(echoServer *echo.Echo) CloseFunc {
	c := jaegertracing.New(echoServer, nil)
	return c.Close
}
