package models

import (
	"database/sql"
	"log"
	"testing"
)

var DB *sql.DB

func ConnectDatabase() error {
	var dbPath string
	if testing.Testing() {
		log.Printf("TESTING---------")
		dbPath = "/Users/krist/code/shop_go/test.db"
	} else {
		dbPath = "/Users/krist/code/shop_go/shop_test.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	DB = db
	//log.Printf("Connecting to DB")
	//log.Printf("DB %+v", db)
	return nil
}
