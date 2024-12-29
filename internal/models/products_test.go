package models

import (
	"fmt"
	"math/rand/v2"
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
	r := rand.IntN(100)
	sku := fmt.Sprintf("FRISK %d", r)
	newProduct := Product{
		Sku:         sku,
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
		Name:        "Something Dice",
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
		Sku:         "GAMELA",
		Name:        "Desk toy",
		Description: "",
		CategoryId:  4,
		Price:       35000,
		Status:      InStock,
	}

	id1, err := AddProductWithQty(instockProduct, 20)
	require.NoError(t, err)
	id2, err := AddProductWithQty(lowstockProduct, 9)
	require.NoError(t, err)
	id3, err := AddProductWithQty(outofstockProduct, -1)
	require.NoError(t, err)

	inv1, _ := GetProductInventory(id1)
	inv2, _ := GetProductInventory(id2)
	inv3, _ := GetProductInventory(id3)

	assert.Equal(t, 20, inv1.Qty)
	assert.Equal(t, 9, inv2.Qty)
	assert.Equal(t, -1, inv3.Qty)

	prod1, _ := GetProductById(id1)
	prod2, _ := GetProductById(id2)
	prod3, _ := GetProductById(id3)

	stat1 := getProductStatus(&prod1)
	stat2 := getProductStatus(&prod2)
	stat3 := getProductStatus(&prod3)

	assert.Equal(t, InStock, stat1)
	assert.Equal(t, LowStock, stat2)
	assert.Equal(t, OutofStock, stat3)
}

func TestSetPreorder(t *testing.T) {
	preorderProduct := Product{
		Sku:         "PRE",
		Name:        "Calendar",
		Description: "",
		CategoryId:  4,
		Price:       35000,
		Status:      OutofStock,
	}

	id1, err := AddProductWithQty(preorderProduct, 0)
	require.NoError(t, err)
	prod1, _ := GetProductById(id1)
	ok, err := SetPreorder(&prod1)
	require.NoError(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, Preorder, prod1.Status)

	//Preorder status can only be applied to products with 0 inventory
	_, err = UpdateProductInventory(id1, 1)
	// updating product inventory does not update product status
	require.NoError(t, err)
	status := getProductStatus(&prod1)
	notok, _ := SetPreorder(&prod1)
	assert.Equal(t, false, notok)
	assert.Equal(t, LowStock, status)
}
