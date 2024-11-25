package models

import (
	"fmt"
	"strings"
)

type Path string

type Photo struct {
	ProductId int    `json:"product_id"`
	Paths     string `json:"images"`
}

func AddPhoto(sku string, filenames []string) error {
	product, err := GetProductBySku(sku)
	if err != nil {
		return err
	}

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
	files := strings.Join(filenames, ", ")
	image := fmt.Sprintf("'[%s]'", files)

	_, err = stmt.Exec(photo.ProductId, image)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func GetPhotosBySku(sku string) (Photo, error) {
	product, err := GetProductBySku(sku)
	if err != nil {
		return Photo{}, err
	}

	stmt, err := DB.Prepare("SELECT product_id, images FROM product_gallery WHERE product_id = ?")
	if err != nil {
		return Photo{}, err
	}
	defer stmt.Close()

	photo := Photo{}
	sqlErr := stmt.QueryRow(product.Id).Scan(&photo.ProductId, &photo.Paths)
	if sqlErr != nil {
		return Photo{}, err
	}

	return photo, nil
}

func GetPhotosById(id int) (Photo, error) {
	stmt, err := DB.Prepare("SELECT product_id, images FROM product_gallery WHERE product_id = ?")
	if err != nil {
		return Photo{}, err
	}
	defer stmt.Close()

	photo := Photo{}
	sqlErr := stmt.QueryRow(id).Scan(&photo.ProductId, &photo.Paths)
	if sqlErr != nil {
		return Photo{}, err
	}

	return photo, nil
}
