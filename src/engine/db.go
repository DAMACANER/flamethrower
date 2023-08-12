package engine

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var max_connections = 100
var DB *sql.DB

func InitDB(dbLoc *string) {
	var err error
	DB, err = sql.Open("sqlite3", *dbLoc)
	if err != nil {
		log.Fatal(err)
	}
	DB.SetMaxOpenConns(max_connections)
}
