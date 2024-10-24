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

/*func TestAddCustomerMissingData(t *testing.T) {
	success, err := AddCustomer(Customer{
		FirstName: "Bob",
		LastName:  "Wood",
		Email:     "bob@wo.od",
	})

	require.NoError(t, err)
	assert.Equal(t, false, success)
}*/
