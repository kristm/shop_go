package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetCategories() ([]Category, error) {
	rows, err := DB.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	categories := make([]Category, 0)

	for rows.Next() {
		category := Category{}
		err = rows.Scan(&category.Id, &category.Name)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return categories, err
}

func GetCategoryById(id int) (Category, error) {
	stmt, err := DB.Prepare("SELECT id, name FROM categories WHERE ID = ?")
	if err != nil {
		return Category{}, err
	}

	category := Category{}
	sqlErr := stmt.QueryRow(id).Scan(&category.Id, &category.Name)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Category{}, nil
		}
		return Category{}, sqlErr
	}
	return category, nil
}