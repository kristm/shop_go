package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddVoucher(t *testing.T) {
	err := AddVoucher(&Voucher{
		TypeId: 3,
		Code:   "NOSHIP",
		Valid:  true,
	})
	assert.Equal(t, err, nil)
}
