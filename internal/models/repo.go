package models

import (
	"database/sql"
	"fmt"
	"log"
	"shop_go/internal/config"
)

var DB *sql.DB
var dbPath string
var dbParams string

func ConnectDatabase(cfg *config.Config) error {
	dbParams = "_foreign_keys=true"
	dbPath = fmt.Sprintf("%s?%s", cfg.DB, dbParams)
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
	dbParams = "" // we don't care about foreign contraints on teardown
	dbPath = fmt.Sprintf("%s?%s", cfg.TEST_DB, dbParams)
	log.Printf("TEST DB path %s\n", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	//db.SetMaxOpenConns(1)
	if err != nil {
		log.Fatal("cannot load DB ", dbPath)
		return err
	}
	DB = db

	return nil
}
