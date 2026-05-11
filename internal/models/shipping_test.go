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
		Fee:        0,
	})

	require.NoError(t, err)
	assert.Greater(t, shippingId, 0)
}

func TestGetShippingById(t *testing.T) {
	shippingId, err := AddShipping(&Shipping{
		CustomerId: 1,
		Status:     0,
		Address:    "Malugay St.",
		City:       "Makati",
		Country:    "PH",
		Zip:        "1203",
		Phone:      "8888",
		Notes:      "X marks the spot",
		Fee:        0,
	})

	shipping, err := GetShippingById(shippingId)
	require.NoError(t, err)
	assert.Equal(t, shipping.Address, "Malugay St.")
	assert.Equal(t, shipping.City, "Makati")
	assert.Equal(t, shipping.Country, "PH")
	assert.Equal(t, shipping.Zip, "1203")
	assert.Equal(t, shipping.Phone, "8888")
	assert.Equal(t, shipping.Notes, "X marks the spot")
}

func TestUpdateShippingStatus(t *testing.T) {
	shippingId, err := AddShipping(&Shipping{
		CustomerId: 1,
		Status:     0,
		Address:    "Malugay St.",
		City:       "Makati",
		Country:    "PH",
		Zip:        "1203",
		Phone:      "8888",
		Notes:      "X marks the spot",
		Fee:        90,
	})

	shipping, err := GetShippingById(shippingId)
	require.NoError(t, err)
	ok, _ := UpdateShippingStatus(shippingId, InTransit)
	shipping, _ = GetShippingById(shippingId)
	assert.Equal(t, ok, true)
	assert.Equal(t, shipping.Status, InTransit)

	ok, _ = UpdateShippingStatus(shippingId, Shipped)
	shipping, _ = GetShippingById(shippingId)
	assert.Equal(t, ok, true)
	assert.Equal(t, shipping.Status, Shipped)

	ok, _ = UpdateShippingStatus(shippingId, Delivered)
	shipping, _ = GetShippingById(shippingId)
	assert.Equal(t, ok, true)
	assert.Equal(t, shipping.Status, Delivered)
}
