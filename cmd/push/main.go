package main

import (
	"fmt"
	"log"
	"os"
	"webup/push"
	"webup/push/http"
	"webup/push/repository/memory"
	"webup/push/repository/mysql"
	"webup/push/repository/postgres"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo/engine/standard"
)

func main() {

	config := push.RuntimeConfig{}
	if _, err := toml.DecodeFile(os.Getenv("CONFIG_FILEPATH"), &config); err != nil {
		log.Fatalln("Unable to parse config: ", err)
	}

	// prepare Token Repository according to config
	var repo push.TokenRepository
	if config.StorageDriver == push.MySQLStorageDriver {
		repo = &mysql.TokenRepository{
			Config: config,
		}
		// clean MySQL connection
		defer mysql.Close()
	} else if config.StorageDriver == push.PostgresStorageDriver {
		repo = &postgres.TokenRepository{
			Config: config,
		}
		// clean Postgres connection
		defer postgres.Close()
	} else {
		repo = new(memory.TokenRepository)
	}

	e := http.SetupRoutes(config, repo)

	fmt.Println("Listening on :3000")
	e.Run(standard.New(":3000"))
}
