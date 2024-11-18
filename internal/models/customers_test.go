package models

import (
	"testing"

	//"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddCustomer(t *testing.T) {
	newCustomer := Customer{
		FirstName: "Bob",
		LastName:  "Wood",
		Email:     "bob@wo.od",
		Phone:     "123-456-789",
	}
	customerId, err := AddCustomer(&newCustomer)

	require.NoError(t, err)
	assert.Greater(t, customerId, 0)
}

func TestAddCustomerMissingData(t *testing.T) {
	success, err := AddCustomer(&Customer{
		FirstName: "Bob",
		LastName:  "Wood",
		Email:     "bob@wo.od",
	})

	assert.EqualError(t, err, "Invalid Customer Data")
	assert.Equal(t, -1, success)
}

func TestDuplicateCustomer(t *testing.T) {
	customer1, err := AddCustomer(&Customer{
		FirstName: "Bob",
		LastName:  "West",
		Email:     "bob@we.st",
		Phone:     "12345",
	})

	duplicateId, err := AddOrGetCustomer(&Customer{
		FirstName: "Bob",
		LastName:  "West",
		Email:     "bob@we.st",
		Phone:     "12345",
	})

	require.NoError(t, err)
	// returns existing id
	assert.Equal(t, customer1, duplicateId)
}
