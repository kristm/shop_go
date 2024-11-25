package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPhoto(t *testing.T) {
	newProduct := Product{
		Sku:        "NWE",
		Name:       "Smile",
		CategoryId: 4,
		Price:      100,
	}
	_, _ = AddProduct(newProduct)
	path := "photo_nwe.jpg"
	err := AddPhoto("NWE", path)
	assert.Equal(t, err, nil)
}
