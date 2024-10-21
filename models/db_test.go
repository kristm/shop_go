package models

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	db, err := sql.Open("sqlite3", "file:../test.db?cache=shared")
	if err != nil {
		return -1, fmt.Errorf("could not connect to database: %w", err)
	}

	defer func() {
		for _, t := range []string{"orders", "customers"} {
			_, _ = db.Exec(fmt.Sprintf("DELETE FROM %s", t))
		}

		db.Close()
	}()

	return m.Run(), nil
}
