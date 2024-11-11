package models

import (
	"database/sql"
	"fmt"
	"log"
	"shop_go/internal/config"
	"testing"
)

var DB *sql.DB

func ConnectDatabase(cfg *config.Config) error {
	var dbPath string
	var err error
	dbParams := "_foreign_keys=true"

	if testing.Testing() {
		dbPath = fmt.Sprintf("%s?%s", cfg.TEST_DB, dbParams)
	} else {
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

func ConnectTestDatabase(cfg *config.Config) error {
	var dbPath string
	var err error
	dbParams := "_foreign_keys=true"

	dbPath = fmt.Sprintf("%s?%s", cfg.TEST_DB, dbParams)

	log.Printf("TEST DB path %s\n", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("cannot load DB ", dbPath)
		return err
	}
	DB = db

	return nil
}
