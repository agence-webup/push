package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"webup/push"

	// MySQL Driver
	_ "github.com/go-sql-driver/mysql"
)

var _db *sql.DB

// GetDB initializes a MySQL connection
func GetDB(config push.MySQLConfig) *sql.DB {
	if _db == nil {
		connectionString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", config.Username, config.Password, config.Hostname, config.Port, config.Database)

		database, err := sql.Open("mysql", connectionString)
		if err != nil {
			log.Println("MySQL init error: ", err)
			return nil
		}
		_db = database

		creationQuery := `CREATE TABLE IF NOT EXISTS push_tokens (
        id int(11) unsigned NOT NULL AUTO_INCREMENT,
        uuid varchar(256) NOT NULL,
        value varchar(1024) NOT NULL DEFAULT '',
        platform tinyint(4) NOT NULL,
        language varchar(6) NOT NULL DEFAULT '',
        created_at datetime NOT NULL,
        PRIMARY KEY (id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
		_, err = _db.Exec(creationQuery)
		if err != nil {
			log.Println("MySQL table creation error: ", err)
			return nil
		}
	}

	err := _db.Ping()
	if err != nil {
		log.Println("MySQL ping error: ", err)
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
