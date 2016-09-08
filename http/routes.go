package http

import (
	"net/http"

	"github.com/labstack/echo"
)

// SetupRoutes prepares the HTTP server for REST API
func SetupRoutes() *echo.Echo {
	e := echo.New()

	e.POST("/send", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	tokenResource := TokenResource{}
	e.POST("/tokens", tokenResource.AddToken())
	e.DELETE("/tokens/:token", tokenResource.RemoveToken())

	return e
}
