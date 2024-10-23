package models

import (
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

func AddCustomer(newCustomer Customer) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	if !ValidateNotEmpty(newCustomer) {
		return false, nil
	}

	stmt, err := tx.Prepare("INSERT INTO customers (first_name, last_name, email, phone) VALUES (?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newCustomer.FirstName, newCustomer.LastName, newCustomer.Email, newCustomer.Phone)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
