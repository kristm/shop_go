package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
)

type Subscriber struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func AddSubscriber(email string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO subscribers (email) VALUES(?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.Exec(email); err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) {
			if sqlErr.Code == sqlite3.ErrConstraint {
				return nil
			}

			return err
		}
	}

	return tx.Commit()
}

func Unsubscribe(email string) error {
	var id int
	err := DB.QueryRow("SELECT id FROM subscribers WHERE email LIKE ?", email).Scan(&id)

	if err != nil {
		return err
	}

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	t := time.Now()
	time := t.Format(DATE_FORMAT)
	newEmail := fmt.Sprintf("%s-%d", email, id)
	if _, err := tx.Exec("UPDATE subscribers SET email = ?, updated_at = ? WHERE id = ?", newEmail, time, id); err != nil {
		return err
	}

	return tx.Commit()
}
