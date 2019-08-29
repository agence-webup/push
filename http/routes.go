package http

import (
	"webup/push"

	"github.com/labstack/echo"
)

// SetupRoutes prepares the HTTP server for REST API
func SetupRoutes(config push.RuntimeConfig, repo push.TokenRepository) *echo.Echo {
	e := echo.New()

	// e.Use(middleware.Logger())

	// -> /send handler
	sendResource := SendResource{
		Config:          config,
		TokenRepository: repo,
	}
	e.POST("/send", sendResource.Send())

	// -> /tokens handlers
	tokenResource := TokenResource{
		Repository: repo,
	}
	e.POST("/tokens", tokenResource.AddToken())
	e.DELETE("/tokens/:platform/:value", tokenResource.RemoveToken())

	return e
}
