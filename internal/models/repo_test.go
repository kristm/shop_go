package models

import (
	"fmt"
	"log"
	"os"
	"shop_go/internal/config"
	"testing"
)

var testTables = []string{"categories", "customers", "orders", "order_products", "shipping", "vouchers"}

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

func ClearProductTestData() (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	query := "DELETE FROM products WHERE sku LIKE ?"
	stmt, err := DB.Prepare(query)
	if err != nil {
		return false, err
	}

	defer stmt.Close()
	_, err = stmt.Exec("WKW%")
	if err != nil {
		return false, err
	}
	tx.Commit()
	log.Printf("product test data cleared")
	return true, nil
}

func TestMain(m *testing.M) {
	log.Println("Test Models Main")
	cfg, err := config.LoadConfig("../../.env")
	if err != nil {
		log.Printf("CFG TEST MAIN %v", &cfg)
		//panic(err)
	}
	log.Printf("OK CFG TEST MAIN %v", &cfg)
	ConnectTestDatabase(&cfg)
	log.Println("Prepare Test tables")
	for _, table := range testTables {
		_, err := ClearTestTable(table)
		if err != nil {
			log.Printf("Teardown error %v", err)
		}
	}
	code := m.Run()

	log.Println("Models Teardown")
	for _, table := range testTables {
		_, err := ClearTestTable(table)
		if err != nil {
			log.Printf("Teardown error %v", err)
		}
	}

	//_, err = ClearProductTestData()
	//if err != nil {
	//	log.Printf("product test teardown error %v", err)
	//}
	os.Exit(code)
}
