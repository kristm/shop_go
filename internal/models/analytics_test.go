package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddCartAnalytics(t *testing.T) {
	analytics := Analytics{
		IpAddress: "127.0.0.1",
		Device:    "Laptop",
		Others:    "cart_age=10.5",
	}

	ok, err := AddCartAnalytics(&analytics)
	require.NoError(t, err)
	assert.Equal(t, true, ok)
}
