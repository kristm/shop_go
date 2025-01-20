package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddVoucher(t *testing.T) {
	now := time.Now()
	expiry := now.AddDate(0, 1, 0)
	err := AddVoucher(&Voucher{
		TypeId:  3,
		Code:    "NOSHIP",
		Valid:   true,
		Expires: expiry.Format(time.RFC3339),
	})
	assert.Equal(t, err, nil)
}

func TestValidateVoucher(t *testing.T) {
	now := time.Now()
	expired := now.AddDate(0, -1, 0)
	validMonth := now.AddDate(0, 1, 0)
	err := AddVoucher(&Voucher{
		TypeId:  2,
		Code:    "EXPIRED",
		Valid:   true,
		Expires: expired.Format(time.RFC3339),
	})
	err = AddVoucher(&Voucher{
		TypeId:  2,
		Code:    "NOW",
		Valid:   true,
		Expires: now.Format(time.RFC3339),
	})
	err = AddVoucher(&Voucher{
		TypeId:  3,
		Code:    "FREESHIP",
		Valid:   true,
		Expires: validMonth.Format(time.RFC3339),
	})

	validVoucher, err := ValidateVoucher("FREESHIP")
	require.NoError(t, err)
	invalidVoucher, err := ValidateVoucher("EXPIRED")
	require.NoError(t, err)
	nowVoucher, err := ValidateVoucher("NOW")
	require.NoError(t, err)

	assert.Equal(t, true, validVoucher)
	assert.Equal(t, false, invalidVoucher)
	assert.Equal(t, true, nowVoucher)
}

func TestApplyVoucher(t *testing.T) {
	now := time.Now()
	validMonth := now.AddDate(0, 1, 0)
	_ = AddVoucher(&Voucher{
		TypeId:  2,
		Code:    "50OFF",
		Valid:   true,
		Expires: validMonth.Format(time.RFC3339),
	})
	price := 165000.00
	err := ApplyVoucher("50OFF", &price)
	require.NoError(t, err)
	assert.Equal(t, 82500.00, price)
}
