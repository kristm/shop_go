package models

import (
	"database/sql"
	"log"
)

type Voucher struct {
	Id      int    `json:"id"`
	TypeId  int    `json:"voucher_type_id"`
	Code    string `json:"code"`
	Valid   bool   `json:"valid"`
	Expires string `json:"expires_at"`
}

func AddVoucher(voucher *Voucher) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO vouchers (voucher_type_id, code, valid, expires_at) VALUES (?,?,?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(voucher.TypeId, voucher.Code, voucher.Valid, voucher.Expires)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil

}

func GetVoucherByCode(code string) (*Voucher, error) {
	stmt, err := DB.Prepare("SELECT id, voucher_type_id, code, valid FROM vouchers WHERE code = ?")
	if err != nil {
		return nil, err
	}

	voucher := Voucher{}
	sqlErr := stmt.QueryRow(code).Scan(&voucher.Id, &voucher.TypeId, &voucher.Code, &voucher.Valid)
	if sqlErr != nil {
		return nil, sqlErr
	}

	return &voucher, nil
}

func ValidateVoucher(code string) (bool, error) {
	stmt, err := DB.Prepare("SELECT id, voucher_type_id, code, valid, expires_at FROM vouchers WHERE code = ? AND valid = TRUE AND datetime('now') < expires_at")
	if err != nil {
		log.Printf("err %v", err)
		return false, err
	}
	defer stmt.Close()

	voucher := Voucher{}
	sqlErr := stmt.QueryRow(code).Scan(&voucher.Id, &voucher.TypeId, &voucher.Code, &voucher.Valid, &voucher.Expires)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return false, nil
		}
		log.Printf("sqlerr %v", sqlErr)
		return false, sqlErr
	}

	return true, nil
}
