package models

import (
	"database/sql"
	"log"
)

const LESS30 = 1
const LESS50 = 2
const FREESHIPPING = 3

type Voucher struct {
	Id             int    `json:"id"`
	TypeId         int    `json:"voucher_type_id"`
	Code           string `json:"code"`
	Valid          bool   `json:"valid"`
	RequiredAmount int    `json:"minimum_spend"`
	Amount         int    `json:"amount,omitempty"`
	Expires        string `json:"expires_at,omitempty"`
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
	stmt, err := DB.Prepare("SELECT v.id, v.voucher_type_id, v.code, v.valid, v.minimum_spend, vt.amount FROM vouchers v LEFT JOIN voucher_types vt ON v.voucher_type_id = vt.id WHERE code = ?")
	if err != nil {
		return nil, err
	}

	voucher := Voucher{}
	sqlErr := stmt.QueryRow(code).Scan(&voucher.Id, &voucher.TypeId, &voucher.Code, &voucher.Valid, &voucher.RequiredAmount, &voucher.Amount)
	if sqlErr != nil {
		return nil, sqlErr
	}

	return &voucher, nil
}

func ValidateVoucher(code string) (bool, error) {
	stmt, err := DB.Prepare("SELECT id, voucher_type_id, code, valid, minimum_spend, expires_at FROM vouchers WHERE code = ? AND valid = TRUE AND datetime('now') < expires_at")
	if err != nil {
		log.Printf("err %v", err)
		return false, err
	}
	defer stmt.Close()

	voucher := Voucher{}
	sqlErr := stmt.QueryRow(code).Scan(&voucher.Id, &voucher.TypeId, &voucher.Code, &voucher.Valid, &voucher.RequiredAmount, &voucher.Expires)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return false, nil
		}
		log.Printf("sqlerr %v", sqlErr)
		return false, sqlErr
	}

	return true, nil
}

func ApplyVoucher(code string, price *float64) error {
	stmt, err := DB.Prepare("SELECT v.voucher_type_id, v.minimum_spend, vt.amount FROM vouchers v LEFT JOIN voucher_types vt ON v.voucher_type_id = vt.id WHERE v.code = ?")

	if err != nil {
		log.Printf("err %v", err)
		return err
	}

	var amount int
	var minimum_spend int
	var voucher_type_id int
	sqlErr := stmt.QueryRow(code).Scan(&voucher_type_id, &minimum_spend, &amount)
	if sqlErr != nil {
		return sqlErr
	}

	if voucher_type_id < FREESHIPPING && int(*price) >= minimum_spend {
		*price = *price - (*price * (float64(amount) * 0.01))
	}

	return nil
}
