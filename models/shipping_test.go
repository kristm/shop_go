package models

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestAddShipping(t *testing.T) {
	success, err := AddShipping(Shipping{
		CustomerId: 1,
		Status:     0,
		Address:    "Malugay St.",
		City:       "Makati",
		Country:    "PH",
		Zip:        "1203",
		Phone:      "8888",
	})

	require.NoError(t, err)
	assert.Equal(t, true, success)
}
