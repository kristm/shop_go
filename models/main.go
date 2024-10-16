package models

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "/Users/krist/code/shop_go/shop_test.db")
	if err != nil {
		return err
	}
	DB = db
	log.Printf("Connecting to DB")
	log.Printf("DB %+v", db)
	return nil
}
