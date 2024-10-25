package models

type Shipping struct {
	Id         int    `json:"id"`
	CustomerId int    `json:"customer_id"`
	Status     int    `json:"status"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Zip        string `json:"zip"`
	Phone      string `json:"phone"`
}

func AddShipping(newAddress Shipping) (int, error) {
	tx, err := DB.Begin()
	if err != nil {
		return -1, err
	}

	//if !ValidateNotEmpty(newAddress) {
	//	return -1, nil
	//}

	stmt, err := tx.Prepare("INSERT INTO shipping (customer_id, status, address, city, country, zip, phone) VALUES (?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return -1, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(newAddress.CustomerId, newAddress.Status, newAddress.Address, newAddress.City, newAddress.Country, newAddress.Zip, newAddress.Phone)

	if err != nil {
		return -1, err
	}

	tx.Commit()
	id, _ := res.LastInsertId()

	return int(id), nil
}
