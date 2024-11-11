package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddShipping(t *testing.T) {
	shippingId, err := AddShipping(&Shipping{
		CustomerId: 1,
		Status:     0,
		Address:    "Malugay St.",
		City:       "Makati",
		Country:    "PH",
		Zip:        "1203",
		Phone:      "8888",
	})

	require.NoError(t, err)
	assert.Greater(t, shippingId, 0)
}
