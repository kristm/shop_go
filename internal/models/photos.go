package models

import "log"

type Path string

type Photo struct {
	ProductId int      `json:"product_id"`
	Paths     []string `json:"images"`
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

	stmt, err := tx.Prepare("INSERT INTO product_gallery (product_id, image) VALUES (?, ?)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	photo := Photo{ProductId: product.Id}

	_, err = stmt.Exec(photo.ProductId, filename)

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

	stmt, err := DB.Prepare("SELECT product_id, image, image2, image3 FROM product_gallery WHERE product_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	photo := Photo{}
	var path1, path2, path3 string
	sqlErr := stmt.QueryRow(product.Id).Scan(&photo.ProductId, &path1, &path2, &path3)
	if sqlErr != nil {
		return nil, err
	}

	photo.Paths = []string{path1, path2, path3}
	return &photo, nil

}
