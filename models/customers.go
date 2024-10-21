package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Customer struct {
	Id        int    `json: "id"`
	FirstName string `json: "first_name"`
	LastName  string `json: "last_name"`
	Email     string `json: "email"`
	Phone     string `json: "phone"`
}

func AddCustomer(newCustomer Customer) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
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
