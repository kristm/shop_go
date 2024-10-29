package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

func TestAddProduct(t *testing.T) {
	timestamp := time.Now().Unix()
	sku := fmt.Sprintf("WKWK %d", timestamp)
	newProduct := Product{
		Sku:         sku,
		Name:        "Something Nice",
		Description: "",
		CategoryId:  4,
		Price:       10000,
		Status:      InStock,
	}

	productId, err := AddProduct(newProduct)
	if err != nil {
		t.Errorf("Error when adding Product %v", err)
	}
	if productId < 0 {
		t.Errorf("Failed to add Product %v", err)
	}

	ok, err := AddProductInventory(productId, 10)
	require.NoError(t, err)
	assert.Equal(t, ok, true)

	inventory, err := GetProductInventory(productId)
	require.NoError(t, err)
	assert.Equal(t, 10, inventory.Qty)
}
