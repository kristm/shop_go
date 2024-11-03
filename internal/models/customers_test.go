package models

import (
	"log"
	"testing"

	//"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(tb testing.TB) func(tb testing.TB) {
	log.Println("setup test")

	return func(tb testing.TB) {
		ClearTestTable("customers")
		log.Println("teardown test")
	}
}

func TestAddCustomer(t *testing.T) {
	customerId, err := AddCustomer(Customer{
		FirstName: "Bob",
		LastName:  "Wood",
		Email:     "bob@wo.od",
		Phone:     "123-456-789",
	})

	require.NoError(t, err)
	assert.Greater(t, customerId, 0)
}

func TestAddCustomerMissingData(t *testing.T) {
	success, err := AddCustomer(Customer{
		FirstName: "Bob",
		LastName:  "Wood",
		Email:     "bob@wo.od",
	})

	require.NoError(t, err)
	assert.Equal(t, -1, success)
}

func TestDuplicateCustomer(t *testing.T) {
	customer1, err := AddCustomer(Customer{
		FirstName: "Bob",
		LastName:  "West",
		Email:     "bob@we.st",
		Phone:     "12345",
	})

	duplicateId, err := AddOrGetCustomer(Customer{
		FirstName: "Bob",
		LastName:  "West",
		Email:     "bob@we.st",
		Phone:     "12345",
	})

	require.NoError(t, err)
	// returns existing id
	assert.Equal(t, customer1, duplicateId)
}
