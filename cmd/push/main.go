package main

import (
	"fmt"
	"webup/push/http"
	"webup/push/repository/mysql"

	"github.com/labstack/echo/engine/standard"
)

func main() {

	e := http.SetupRoutes()

	// clean MySQL connection
	defer mysql.Close()

	fmt.Println("Listening on http://localhost:3000")
	e.Run(standard.New(":3000"))
}
