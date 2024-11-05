package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddSocials(t *testing.T) {
	ok, err := AddSocials(&Socials{
		CustomerId: 1,
		Subscribe:  true,
		Socials:    "Chik Chok",
	})
	require.NoError(t, err)
	assert.Equal(t, ok, true)
}
