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
	var cfg config.Config
	var err error
	dbParams := "_foreign_keys=true"

	if testing.Testing() {
		cfg, err = config.LoadConfig("../.env")
		if err != nil {
			log.Fatal("cannot load config ", err)
		}

		dbPath = fmt.Sprintf("%s?%s", cfg.TEST_DB, dbParams)
	} else {
		cfg, err = config.LoadConfig(".env")
		if err != nil {
			log.Fatal("cannot load config ", err)
		}
		dbPath = fmt.Sprintf("%s?%s", cfg.DB, dbParams)
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
