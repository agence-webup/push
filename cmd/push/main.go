package main

import (
	"fmt"
	"log"
	"os"
	"webup/push"
	"webup/push/http"
	"webup/push/repository/mysql"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo/engine/standard"
)

func main() {

	config := push.RuntimeConfig{}
	if _, err := toml.DecodeFile(os.Getenv("CONFIG_FILEPATH"), &config); err != nil {
		log.Fatalln("Unable to parse config: ", err)
	}

	e := http.SetupRoutes(config)

	// clean MySQL connection
	defer mysql.Close()

	fmt.Println("Listening on http://localhost:3000")
	e.Run(standard.New(":3000"))
}
