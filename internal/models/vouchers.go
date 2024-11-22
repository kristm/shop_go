package models

type Voucher struct {
	Id     int    `json:"id"`
	TypeId int    `json:"voucher_type_id"`
	Code   string `json:"code"`
	Valid  bool   `json:"valid"`
}

func AddVoucher(voucher *Voucher) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO vouchers (voucher_type_id, code, valid) VALUES (?,?,?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(voucher.TypeId, voucher.Code, voucher.Valid)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil

}
