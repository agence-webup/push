package http

import "github.com/labstack/echo"

// SetupRoutes prepares the HTTP server for REST API
func SetupRoutes() *echo.Echo {
	e := echo.New()

	sendResource := SendResource{}
	e.POST("/send", sendResource.Send())

	tokenResource := TokenResource{}
	e.POST("/tokens", tokenResource.AddToken())
	e.DELETE("/tokens/:token", tokenResource.RemoveToken())

	return e
}
