package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	CategoryId   int    `json:"category_id"`
	PriceInCents int    `json:"price_in_cents"`
}

func GetProducts(category_id int) ([]Product, error) {
	stmt, err := DB.Prepare("SELECT id, name, description, category_id, price_in_cents FROM products WHERE category_id = ?")
	if err != nil {
		return nil, err
	}

	rows, sqlErr := stmt.Query(category_id)
	if sqlErr != nil {
		return nil, sqlErr
	}

	defer rows.Close()
	products := make([]Product, 0)

	for rows.Next() {
		product := Product{}
		err = rows.Scan(&product.Id, &product.Name, &product.Description, &product.CategoryId, &product.PriceInCents)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	err = rows.Err()

	return products, nil
}
