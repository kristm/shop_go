package models

import (
	"fmt"
	"log"
)

type Path string

type Photo struct {
	ProductId int    `json:"product_id"`
	Paths     string `json:"images"`
}

func AddPhoto(sku string, filename string) error {
	product, err := GetProductBySku(sku)
	if err != nil {
		return err
	}

	log.Printf("PRODUCT %v ", product)

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO product_gallery (product_id, images) VALUES (?, ?)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	photo := Photo{ProductId: product.Id}
	image := fmt.Sprintf("'[%s]'", filename)

	_, err = stmt.Exec(photo.ProductId, image)

	if err != nil {
		log.Printf("ERROR FOUND")
		return err
	}

	tx.Commit()

	return nil
}

func GetPhotosBySku(sku string) (*Photo, error) {
	product, err := GetProductBySku(sku)
	if err != nil {
		return nil, err
	}

	stmt, err := DB.Prepare("SELECT product_id, images FROM product_gallery WHERE product_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	photo := Photo{}
	sqlErr := stmt.QueryRow(product.Id).Scan(&photo.ProductId, &photo.Paths)
	if sqlErr != nil {
		return nil, err
	}

	return &photo, nil

}
