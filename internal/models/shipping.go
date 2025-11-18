package models

type ShippingStatus int

const (
	Packing ShippingStatus = iota
	InTransit
	Shipped
	Delivered
)

type Shipping struct {
	Id         int            `json:"id"`
	CustomerId int            `json:"customer_id"`
	Status     ShippingStatus `json:"status"`
	Address    string         `json:"address"`
	City       string         `json:"city"`
	Country    string         `json:"country"`
	Zip        string         `json:"zip"`
	Phone      string         `json:"phone"`
	Notes      string         `json:"notes"`
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

func GetShippingById(id int) (*Shipping, error) {
	stmt, err := DB.Prepare("SELECT status, address, city, country, zip, phone, notes FROM shipping WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	shipping := &Shipping{}
	sqlErr := stmt.QueryRow(id).Scan(&shipping.Status, &shipping.Address, &shipping.City, &shipping.Country, &shipping.Zip, &shipping.Phone, &shipping.Notes)
	if sqlErr != nil {
		return nil, sqlErr
	}

	return shipping, nil
}

func UpdateShippingStatus(id int, status ShippingStatus) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE shipping SET STATUS = ? WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(status, id)
	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
