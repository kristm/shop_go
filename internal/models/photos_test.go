package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddPhoto(t *testing.T) {
	newProduct := Product{
		Sku:        "NWE",
		Name:       "Smile",
		CategoryId: 4,
		Price:      100,
	}
	_, _ = AddProduct(newProduct)
	path := []string{"photo_nwe.jpg"}
	err := AddPhoto("NWE", path)
	assert.Equal(t, err, nil)
}

func TestAddMultiplePhotos(t *testing.T) {
	newProduct := Product{
		Sku:        "NWE2",
		Name:       "Smirk",
		CategoryId: 4,
		Price:      200,
	}
	_, _ = AddProduct(newProduct)
	path := []string{"photo_smirk.jpg", "photo_smirk2.jpg"}
	err := AddPhoto("NWE2", path)
	assert.Equal(t, err, nil)
}

func TestGetPhotosBySku(t *testing.T) {
	newProduct := Product{
		Sku:        "NWE3",
		Name:       "Cheese",
		CategoryId: 4,
		Price:      150,
	}
	_, _ = AddProduct(newProduct)
	path := []string{"photo_cheese.jpg", "photo_butter.jpg"}
	err := AddPhoto("NWE3", path)

	photo, err := GetPhotosBySku("NWE3")
	require.NoError(t, err)
	assert.Equal(t, photo.Paths, "'[photo_cheese.jpg, photo_butter.jpg]'")
}
