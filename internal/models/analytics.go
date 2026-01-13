package models

type Analytics struct {
	CustomerId int    `json:"customer_id"`
	IpAddress  string `json:"ip_address"`
	Device     string `json:"device"`
	Others     string `json:"others"` //CartAge is in minutes
}

func AddAnalytics(analytics *Analytics) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO analytics (customer_id, ip_address, device, others) VALUES (?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(analytics.CustomerId, analytics.IpAddress, analytics.Device, analytics.Others)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func AddCartAnalytics(analytics *Analytics) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO analytics (ip_address, device, others) VALUES (?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(analytics.IpAddress, analytics.Device, analytics.Others)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
