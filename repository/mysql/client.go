package mysql

import (
	"database/sql"
	"log"

	// MySQL Driver
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// GetDB initializes a MySQL connection
func GetDB() *sql.DB {
	if db == nil {
		database, err := sql.Open("mysql", "root:pass@tcp(localhost:3306)/db")
		if err != nil {
			log.Fatalln("MySQL init error: ", err)
			return nil
		}
		db = database
	}

	err := db.Ping()
	if err != nil {
		log.Fatalln("MySQL ping error: ", err)
		return nil
	}

	return db
}

// Close cleans the MySQL connection
func Close() {
	db.Close()
	db = nil
}
