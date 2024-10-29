package models

import (
	"database/sql"
	"fmt"
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

func ClearTestTable(tableName string) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	sql := fmt.Sprintf("DELETE FROM %s WHERE id >= ?", tableName)
	stmt, err := DB.Prepare(sql)
	if err != nil {
		return false, err
	}

	defer stmt.Close()
	_, err = stmt.Exec(1)
	if err != nil {
		return false, err
	}
	tx.Commit()
	log.Printf("%s table cleared", tableName)
	return true, nil
}

func TestMain(m *testing.M) {
	log.Println("Test Models Main")
	ConnectTestDatabase()
	code := m.Run()

	testTables := []string{"customers", "products", "product_inventory", "orders", "order_products", "shipping"}
	log.Println("Models Teardown")
	for _, table := range testTables {
		_, err := ClearTestTable(table)
		if err != nil {
			log.Printf("Teardown error %v", err)
		}
	}
	os.Exit(code)
}
