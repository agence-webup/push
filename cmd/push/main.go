package main

import (
	"fmt"
	"webup/push/http"

	"github.com/labstack/echo/engine/standard"
)

func main() {

	e := http.SetupRoutes()

	fmt.Println("Listening on http://localhost:3000")
	e.Run(standard.New(":3000"))
}
