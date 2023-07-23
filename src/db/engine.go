package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbLocation string) {
	var err error
	DB, err = sql.Open("sqlite3", dbLocation)
	if err != nil {
		log.Fatal(err)
	}

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxIdleTime(time.Minute * 5)
}
