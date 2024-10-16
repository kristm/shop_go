package models

import (
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

// wonky
//func (prod *Product) MarshalJSON() ([]byte, error) {
//
//	log.Printf("%+v", prod)
//	prod.Price = prod.Price / 100
//	var out string
//	out = fmt.Sprintf("%v", prod)
//	return []byte(out[1:len(out)]), nil
//}

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

		//need to convert to precision float
		product.Price = float64(int(product.Price)) / 100.00

		products = append(products, product)
	}
	err = rows.Err()

	return products, nil
}
