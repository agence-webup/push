package http

import (
	"webup/push"
	"webup/push/repository/memory"
	"webup/push/repository/mysql"

	"github.com/labstack/echo"
)

// SetupRoutes prepares the HTTP server for REST API
func SetupRoutes(config push.RuntimeConfig) *echo.Echo {
	e := echo.New()

	// e.Use(middleware.Logger())

	// prepare Token Repository according to config
	var repo push.TokenRepository
	if config.StorageDriver == push.MySQLStorageDriver {
		repo = &mysql.TokenRepository{
			Config: config,
		}
	} else {
		repo = new(memory.TokenRepository)
	}

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
