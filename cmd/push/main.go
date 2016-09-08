package main

import (
	"webup/push/http"

	"github.com/labstack/echo/engine/standard"
)

func main() {

	e := http.SetupRoutes()

	e.Run(standard.New(":3000"))
}
