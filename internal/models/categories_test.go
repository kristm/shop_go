package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCategories(t *testing.T) {
	category1 := Category{Name: "Music", Enabled: true}
	category2 := Category{Name: "Cinema", Enabled: false}
	category3 := Category{Name: "Lit", Enabled: true}
	_ = AddCategory(&category1)
	_ = AddCategory(&category2)
	_ = AddCategory(&category3)

	categories, err := GetCategories()

	require.NoError(t, err)
	assert.Equal(t, 2, len(categories))
}
