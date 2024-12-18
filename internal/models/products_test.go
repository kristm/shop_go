package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
}

func TestGetProductBySku(t *testing.T) {
	newProduct := Product{
		Sku:         "FISKBO",
		Name:        "Frame",
		Description: "",
		CategoryId:  4,
		Price:       15000,
		Status:      InStock,
	}

	productId, err := AddProduct(newProduct)
	require.NoError(t, err)
	assert.Greater(t, productId, 0)
}

func TestAddProductInventory(t *testing.T) {
	timestamp := time.Now().Unix()
	sku := fmt.Sprintf("WKWS %d", timestamp)
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

func TestUpdateProductInventory(t *testing.T) {
	timestamp := time.Now().Unix()
	sku := fmt.Sprintf("WKWZ %d", timestamp)
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

	updated, err := UpdateProductInventory(productId, 20)
	require.NoError(t, err)
	assert.Equal(t, updated, true)

	inventory, err := GetProductInventory(productId)
	require.NoError(t, err)
	assert.Equal(t, 20, inventory.Qty)
}

func TestAddProductWithQty(t *testing.T) {
	timestamp := time.Now().Unix()
	sku := fmt.Sprintf("WKW1 %d", timestamp)
	newProduct := Product{
		Sku:         sku,
		Name:        "Something Rice",
		Description: "",
		CategoryId:  4,
		Price:       10000,
		Status:      InStock,
	}
	productId, err := AddProductWithQty(newProduct, 100)
	if err != nil {
		t.Errorf("Error when adding Product %v", err)
	}
	if productId < 0 {
		t.Errorf("Failed to add Product %v", err)
	}
	inventory, err := GetProductInventory(productId)
	require.NoError(t, err)
	assert.Equal(t, 100, inventory.Qty)
}

func TestGetProductStatus(t *testing.T) {
	instockProduct := Product{
		Sku:         "FISKBO",
		Name:        "Frame",
		Description: "",
		CategoryId:  4,
		Price:       15000,
		Status:      InStock,
	}
	lowstockProduct := Product{
		Sku:         "DURIAN",
		Name:        "Fruit Cup",
		Description: "",
		CategoryId:  4,
		Price:       25000,
		Status:      InStock,
	}
	outofstockProduct := Product{
		Sku:         "ZOID",
		Name:        "Desk toy",
		Description: "",
		CategoryId:  4,
		Price:       35000,
		Status:      InStock,
	}

	id1, _ := AddProductWithQty(instockProduct, 20)
	id2, _ := AddProductWithQty(lowstockProduct, 9)
	id3, _ := AddProductWithQty(outofstockProduct, -1)

	prod1, _ := GetProductById(id1)
	prod2, _ := GetProductById(id2)
	prod3, _ := GetProductById(id3)

	stat1 := getProductStatus(&prod1)
	stat2 := getProductStatus(&prod2)
	stat3 := getProductStatus(&prod3)

	assert.Equal(t, stat1, InStock)
	assert.Equal(t, 0, prod1.Qty)
	assert.Equal(t, stat2, LowStock)
	assert.Equal(t, stat3, OutofStock)
}
