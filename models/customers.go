package models

import (
	"fmt"
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

func checkCustomerDuplicate(customer Customer) (int, error) {
	sql := fmt.Sprintf("SELECT id, first_name, last_name, email, phone FROM customers WHERE first_name = ? AND last_name = ? AND email = ? AND phone = ? ORDER BY created_at LIMIT 1")
	stmt, err := DB.Prepare(sql)
	if err != nil {
		return -1, err
	}

	existingCustomer := Customer{}
	sqlErr := stmt.QueryRow(customer.FirstName, customer.LastName, customer.Email, customer.Phone).Scan(&existingCustomer.Id, &existingCustomer.FirstName, &existingCustomer.LastName, &existingCustomer.Email, &existingCustomer.Phone)
	if sqlErr != nil {
		return -1, sqlErr
	}
	log.Printf("existing Customer %v", existingCustomer)

	return existingCustomer.Id, nil
}

func AddOrGetCustomer(customer Customer) (int, error) {
	customerId, _ := checkCustomerDuplicate(customer)
	if customerId > 0 {
		return customerId, nil
	}

	return AddCustomer(customer)
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
