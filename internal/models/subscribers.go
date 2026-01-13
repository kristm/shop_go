package models

import (
	"errors"

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

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
