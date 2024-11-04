package models

import (
	"database/sql"
	"fmt"
	"log"
	"shop_go/internal/config"
	"testing"
)

var DB *sql.DB

func ConnectDatabase() error {
	var dbPath string
	dbParams := "_foreign_keys=true"

	config, err := config.LoadConfig("../.env")
	if err != nil {
		log.Fatal("cannot load config ", err)
	}

	if testing.Testing() {
		dbPath = fmt.Sprintf("%s?%s", config.TEST_DB, dbParams)
	} else {
		dbPath = fmt.Sprintf("%s?%s", config.DB, dbParams)
	}

	log.Printf("DB path %s\n", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("cannot load DB ", dbPath)
		return err
	}
	DB = db

	return nil
}
