package http

import (
	"webup/push/repository/mysql"

	"github.com/labstack/echo"
)

// SetupRoutes prepares the HTTP server for REST API
func SetupRoutes() *echo.Echo {
	e := echo.New()

	// e.Use(middleware.Logger())

	mysqlRepo := new(mysql.TokenRepository)

	sendResource := SendResource{
		TokenRepository: mysqlRepo,
	}
	e.POST("/send", sendResource.Send())

	tokenResource := TokenResource{
		Repository: mysqlRepo,
	}
	e.POST("/tokens", tokenResource.AddToken())
	e.DELETE("/tokens/:platform/:value", tokenResource.RemoveToken())

	return e
}
