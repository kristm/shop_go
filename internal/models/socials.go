package models

import "fmt"

type Socials struct {
	Id         int    `json:"id"`
	CustomerId int    `json:"customer_id"`
	Subscribe  bool   `json:"subscribe"`
	Socials    string `json:"socials"`
}

// check if record already exists
func checkSocials(socials *Socials) (bool, error) {
	sqlQuery := fmt.Sprintf("SELECT customer_id, subscribed_to_newsletter, account_url FROM socials WHERE customer_id = ? AND subscribed_to_newsletter = ? AND account_url = ? ORDER BY created_at LIMIT 1")
	stmt, err := DB.Prepare(sqlQuery)
	if err != nil {
		return false, err
	}

	existingSocials := Socials{}
	sqlErr := stmt.QueryRow(socials.CustomerId, socials.Subscribe, socials.Socials).Scan(&existingSocials.Id, &existingSocials.CustomerId, &existingSocials.Subscribe, &existingSocials.Socials)
	if sqlErr != nil {
		return false, sqlErr
	}

	return true, nil
}

func AddCustomerSocials(socials *Socials) (bool, error) {
	exists, err := checkSocials(socials)
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	ok, err := AddSocials(socials)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func AddSocials(socials *Socials) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO socials (customer_id, subscribed_to_newsletter, account_url) VALUES (?,?,?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(socials.CustomerId, socials.Subscribe, socials.Socials)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
