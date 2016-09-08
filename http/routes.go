package http

import (
	"webup/push/repository/memory"

	"github.com/labstack/echo"
)

// SetupRoutes prepares the HTTP server for REST API
func SetupRoutes() *echo.Echo {
	e := echo.New()

	sendResource := SendResource{}
	e.POST("/send", sendResource.Send())

	tokenResource := TokenResource{
		Repository: new(memory.TokenRepository),
	}
	e.POST("/tokens", tokenResource.AddToken())
	e.DELETE("/tokens/:platform/:value", tokenResource.RemoveToken())

	return e
}
