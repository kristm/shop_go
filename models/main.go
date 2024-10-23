package models

import (
	"database/sql"
)

var DB *sql.DB

func ConnectDatabase() error {
	var dbPath string
	dbPath = "/Users/krist/code/shop_go/shop_test.db"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	DB = db
	//log.Printf("Connecting to DB")
	//log.Printf("DB %+v", db)
	return nil
}
