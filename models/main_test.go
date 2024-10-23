package models

import (
	"database/sql"
	"log"
	"os"
	"testing"
)

func ConnectTestDatabase() {
	db, err := sql.Open("sqlite3", "/Users/krist/code/shop_go/test.db")
	if err != nil {
		log.Println(err)
	}
	DB = db
}

func TestMain(m *testing.M) {
	log.Println("Test Models Main")
	ConnectTestDatabase()
	code := m.Run()

	os.Exit(code)
}
