package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddCustomer(t *testing.T) {
	success, err := AddCustomer(Customer{
		FirstName: "Bob",
		LastName:  "Wood",
		Email:     "bob@wo.od",
		Phone:     "123-456-789",
	})

	require.NoError(t, err)
	assert.Equal(t, true, success)
}
