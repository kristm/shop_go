package main

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

var DB *sql.DB

func setupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("setup suite")

	return func(tb testing.TB) {
		log.Println("teardown suite")
	}
}

func setupTest(tb testing.TB) func(tb testing.TB) {
	log.Println("setup test")

	return func(tb testing.TB) {
		log.Println("teardown test")
	}
}

func ConnectTestDatabase() {
	db, err := sql.Open("sqlite3", "/Users/krist/code/shop_go/test.db")
	if err != nil {
		log.Println(err)
	}
	DB = db
}

func ClearProducts() (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE FROM products WHERE id >= ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()
	_, err = stmt.Exec(1)
	if err != nil {
		return false, err
	}
	tx.Commit()
	log.Println("Products table cleared")
	return true, nil
}

func TestMain(m *testing.M) {
	log.Println("Test Main")
	ConnectTestDatabase()
	code := m.Run()
	_, err := ClearProducts()
	if err != nil {
		log.Printf("Teardown error %v\n", err)
	}
	os.Exit(code)
}

func TestNothing(t *testing.T) {
	assert.Equal(t, true, true)
}
