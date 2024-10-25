package models

import (
	"log"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

type Customer struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func PrintCustomer(customer Customer) Customer {
	log.Printf("------>> %v\n", DB)
	tx, err := DB.Begin()
	if err != nil {
		log.Printf("error %v", err)
	}
	tx.Commit()
	return customer
}

func ValidateNotEmpty(customer Customer) bool {
	values := reflect.ValueOf(customer)
	for i := 0; i < values.NumField(); i++ {
		f := values.Field(i)
		value := values.Field(i).Interface()
		valueType := f.Type().String()
		if valueType == "string" && len(value.(string)) == 0 {
			return false
		}

	}
	return true
}

func AddCustomer(newCustomer Customer) (int, error) {
	tx, err := DB.Begin()
	if err != nil {
		return -1, err
	}

	if !ValidateNotEmpty(newCustomer) {
		return -1, nil
	}

	stmt, err := tx.Prepare("INSERT INTO customers (first_name, last_name, email, phone) VALUES (?, ?, ?, ?)")

	if err != nil {
		return -1, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(newCustomer.FirstName, newCustomer.LastName, newCustomer.Email, newCustomer.Phone)

	if err != nil {
		return -1, err
	}

	tx.Commit()
	id, _ := res.LastInsertId()

	return int(id), nil
}
