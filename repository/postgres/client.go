package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"webup/push"

	// Postgres Driver
	_ "github.com/lib/pq"
)

var _db *sql.DB

// GetDB initializes a Postgres connection
func GetDB(config push.PostgresConfig) *sql.DB {
	if _db == nil {
		connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", config.Username, config.Password, config.Hostname, config.Port, config.Database)

		database, err := sql.Open("postgres", connectionString)
		if err != nil {
			log.Println("Postgres init error: ", err)
			return nil
		}
		_db = database

		creationQuery := `CREATE TABLE IF NOT EXISTS ` + getTableName(config) + ` (
        id SERIAL PRIMARY KEY,
        uuid varchar(256) NOT NULL,
        value varchar(1024) NOT NULL DEFAULT '',
        platform int2 NOT NULL,
        language varchar(6) NOT NULL DEFAULT '',
        created_at timestamp NOT NULL
		);`
		_, err = _db.Exec(creationQuery)
		if err != nil {
			log.Println("Postgres table creation error: ", err)
			return nil
		}
	}

	err := _db.Ping()
	if err != nil {
		log.Println("Postgres ping error: ", err)
		_db = nil
		return nil
	}

	return _db
}

// Close cleans the MySQL connection
func Close() {
	_db.Close()
	_db = nil
}

func getTableName(config push.PostgresConfig) string {
	return config.Prefix + `push_tokens`
}
