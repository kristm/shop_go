package models

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	Id          int     `json:"id"`
	Sku         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryId  int     `json:"category_id"`
	Price       float64 `json:"price"`
}

func (prod Product) MarshalJSON() ([]byte, error) {
	type Alias Product
	computedPrice := float64(int(prod.Price)) / 100.00
	return json.Marshal(&struct {
		*Alias
		Price float64 `json:"price"`
	}{
		Alias: (*Alias)(&prod),
		Price: computedPrice,
	})
}

func GetProducts(category_id int) ([]Product, error) {
	stmt, err := DB.Prepare("SELECT id, name, sku, description, category_id, price_in_cents FROM products WHERE category_id = ?")
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
		err = rows.Scan(&product.Id, &product.Name, &product.Sku, &product.Description, &product.CategoryId, &product.Price)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	err = rows.Err()

	return products, nil
}

func GetProductById(id int) (Product, error) {
	stmt, err := DB.Prepare("SELECT id, name, sku, description, category_id, price_in_cents FROM products WHERE id = ?")
	if err != nil {
		return Product{}, err
	}

	product := Product{}
	sqlErr := stmt.QueryRow(id).Scan(&product.Id, &product.Name, &product.Sku, &product.Description, &product.CategoryId, &product.Price)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Product{}, nil
		}
		return Product{}, sqlErr
	}
	return product, nil
}

func GetProductBySku(sku string) (Product, error) {
	stmt, err := DB.Prepare("SELECT id, name, sku, description, category_id, price_in_cents FROM products WHERE sku = ?")
	if err != nil {
		return Product{}, err
	}

	product := Product{}
	sqlErr := stmt.QueryRow(sku).Scan(&product.Id, &product.Name, &product.Sku, &product.Description, &product.CategoryId, &product.Price)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Product{}, nil
		}
		return Product{}, sqlErr
	}
	return product, nil
}

func AddProduct(newProduct Product) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO products (name, sku, description, category_id, price_in_cents) VALUES (?, ?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newProduct.Name, newProduct.Sku, newProduct.Description, newProduct.CategoryId, newProduct.Price)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
