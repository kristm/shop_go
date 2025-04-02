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
	Notes      string `json:"notes"`
}

func AddShipping(newAddress *Shipping) (int, error) {
	tx, err := DB.Begin()
	if err != nil {
		return -1, err
	}

	//if !ValidateNotEmpty(newAddress) {
	//	return -1, nil
	//}

	stmt, err := tx.Prepare("INSERT INTO shipping (customer_id, status, address, city, country, zip, phone,notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return -1, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(newAddress.CustomerId, newAddress.Status, newAddress.Address, newAddress.City, newAddress.Country, newAddress.Zip, newAddress.Phone, newAddress.Notes)

	if err != nil {
		return -1, err
	}

	tx.Commit()
	id, _ := res.LastInsertId()

	return int(id), nil
}

func GetShippingAddresses() ([]Shipping, error) {
	rows, err := DB.Query("SELECT customer_id, address, city, country, zip, phone, notes FROM shipping ORDER BY customer_id")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var addresses []Shipping
	for rows.Next() {
		var address Shipping
		if err := rows.Scan(&address.CustomerId, &address.Address, &address.City, &address.Country, &address.Zip, &address.Phone, &address.Notes); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return addresses, err
}
