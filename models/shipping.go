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

func AddShipping(newAddress Shipping) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	//if !ValidateNotEmpty(newAddress) {
	//	return false, nil
	//}

	stmt, err := tx.Prepare("INSERT INTO shipping (customer_id, status, address, city, country, zip, phone) VALUES (?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newAddress.CustomerId, newAddress.Status, newAddress.Address, newAddress.City, newAddress.Country, newAddress.Zip, newAddress.Phone)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
