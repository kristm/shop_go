package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
