package models

type Path string

type Photo struct {
	ProductId int    `json:"product_id"`
	Paths     []Path `json:"images"`
}

func AddPhoto(sku string, filename Path) error {
	product, err := GetProductBySku(sku)
	if err != nil {
		return err
	}

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO product_gallery (product_id, image, image2, image3) VALUES (?, ?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	photos := []Path{filename}

	photo := Photo{ProductId: product.Id, Paths: photos}
	_, err = stmt.Exec(photo.ProductId, photo.Paths)

	if err != nil {
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

	photo := Photo{}
	var path1, path2, path3 Path
	sqlErr := stmt.QueryRow(product.Id).Scan(&photo.ProductId, &path1, &path2, &path3)
	if sqlErr != nil {
		return nil, err
	}

	photo.Paths = []Path{path1, path2, path3}
	return &photo, nil

}
