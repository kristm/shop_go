package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type ProductStatus int

const (
	InStock ProductStatus = iota
	LowStock
	OutofStock
)

type Product struct {
	Id          int           `json:"id"`
	Sku         string        `json:"sku"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	CategoryId  int           `json:"category_id"`
	Price       float64       `json:"price"`
	Status      ProductStatus `json:"status"`
	Photos      Photo         `json:"images"`
}

type Inventory struct {
	Id        int `json:"id"`
	ProductId int `json:"product_id"`
	Qty       int `json:"qty"`
}

func (prod Product) MarshalJSON() ([]byte, error) {
	type Alias Product
	computedPrice := float64(int(prod.Price)) / 100.00
	photosCollection := strings.Split(prod.Photos.Paths, ", ")
	return json.Marshal(&struct {
		*Alias
		Price  float64  `json:"price"`
		Photos []string `json:"images"`
	}{
		Alias:  (*Alias)(&prod),
		Price:  computedPrice,
		Photos: photosCollection,
	})
}

func GetAllProducts() ([]Product, error) {
	rows, err := DB.Query("SELECT id, name, sku, description, category_id, price_in_cents, status FROM products")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	products := make([]Product, 0)

	for rows.Next() {
		product := Product{}
		err = rows.Scan(&product.Id, &product.Name, &product.Sku, &product.Description, &product.CategoryId, &product.Price, &product.Status)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return products, err
}

func GetProducts(category_id int) ([]Product, error) {
	stmt, err := DB.Prepare("SELECT id, name, sku, description, category_id, price_in_cents FROM products WHERE category_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

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

		photos, _ := GetPhotosById(product.Id)
		product.Photos = photos

		products = append(products, product)
	}
	//TODO handle error
	err = rows.Err()

	return products, nil
}

func GetProductById(id int) (Product, error) {
	stmt, err := DB.Prepare("SELECT id, name, sku, description, category_id, price_in_cents FROM products WHERE id = ?")
	if err != nil {
		return Product{}, err
	}
	defer stmt.Close()

	product := Product{}
	sqlErr := stmt.QueryRow(id).Scan(&product.Id, &product.Name, &product.Sku, &product.Description, &product.CategoryId, &product.Price)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Product{}, nil
		}
		return Product{}, sqlErr
	}

	product.Status = getProductStatus(&product)

	photos, _ := GetPhotosById(product.Id)
	product.Photos = photos

	return product, nil
}

func GetProductBySku(sku string) (Product, error) {
	stmt, err := DB.Prepare("SELECT id, name, sku, description, category_id, price_in_cents FROM products WHERE sku = ?")
	if err != nil {
		return Product{}, err
	}
	defer stmt.Close()

	product := Product{}
	sqlErr := stmt.QueryRow(sku).Scan(&product.Id, &product.Name, &product.Sku, &product.Description, &product.CategoryId, &product.Price)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Product{}, nil
		}
		return Product{}, sqlErr
	}

	product.Status = getProductStatus(&product)

	photos, _ := GetPhotosById(product.Id)
	product.Photos = photos

	return product, nil
}

func getProductStatus(product *Product) ProductStatus {
	var status ProductStatus
	inventory, _ := GetProductInventory(product.Id)
	qty := inventory.Qty

	switch {
	case qty < 0:
		status = OutofStock
	case qty < 10:
		status = LowStock
	default:
		status = InStock
	}

	return status
}

// move this to cli tool
func Validate(product *Product, fieldname string) bool {
	value := strings.ToLower(product.Sku)
	if strings.ToLower(value) == fieldname {
		return false
	}

	return true
}

func AddProductWithQty(newProduct Product, qty int) (int, error) {
	productId, err := AddProduct(newProduct)
	if err != nil {
		return -1, err
	}

	_, err = AddProductInventory(productId, qty)
	if err != nil {
		return -1, err
	}

	return productId, nil
}

func AddProduct(newProduct Product) (int, error) {
	isValid := Validate(&newProduct, "sku")
	if !isValid {
		return -1, nil
	}

	tx, err := DB.Begin()
	if err != nil {
		return -1, err
	}

	stmt, err := tx.Prepare("INSERT INTO products (name, sku, description, category_id, price_in_cents) VALUES (?, ?, ?, ?, ?)")

	if err != nil {
		return -1, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(newProduct.Name, newProduct.Sku, newProduct.Description, newProduct.CategoryId, newProduct.Price)

	if err != nil {
		return -1, err
	}

	tx.Commit()
	id, _ := res.LastInsertId()

	return int(id), nil
}

func GetProductInventory(productId int) (Inventory, error) {
	query := fmt.Sprintf("SELECT id, product_id, qty FROM product_inventory WHERE product_id = ? ORDER BY created_at LIMIT 1")
	stmt, err := DB.Prepare(query)
	if err != nil {
		return Inventory{}, err
	}
	defer stmt.Close()

	inventory := Inventory{}
	sqlErr := stmt.QueryRow(productId).Scan(&inventory.Id, &inventory.ProductId, &inventory.Qty)
	if sqlErr != nil && sqlErr != sql.ErrNoRows {
		return Inventory{}, nil
	}

	return inventory, nil
}

func UpdateProductInventory(productId int, qty int) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE product_inventory SET qty = ? WHERE product_id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(qty, productId)
	if err != nil {
		return false, err
	}

	tx.Commit()
	return true, nil
}

func AddProductInventory(productId int, qty int) (bool, error) {
	inventory, err := GetProductInventory(productId)
	if err != nil {
		return false, err
	}

	if inventory.Id > 0 {
		UpdateProductInventory(productId, qty)
		return true, nil
	} else {
		//add inventory
		tx, err := DB.Begin()
		if err != nil {
			return false, err
		}

		stmt, err := tx.Prepare("INSERT INTO product_inventory (product_id, qty) VALUES (?, ?)")

		if err != nil {
			return false, err
		}

		defer stmt.Close()
		inventory = Inventory{ProductId: productId, Qty: qty}

		_, err = stmt.Exec(inventory.ProductId, inventory.Qty)

		if err != nil {
			return false, err
		}

		tx.Commit()
		return true, nil
	}

}
