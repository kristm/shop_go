package models

import (
	"testing"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

func TestAddProduct(t *testing.T) {
	newProduct := Product{
		Sku:         "wkwk",
		Name:        "Something Nice",
		Description: "",
		CategoryId:  4,
		Price:       10000,
	}

	got, err := AddProduct(newProduct)
	if err != nil {
		t.Errorf("Error when adding Product %v", err)
	}
	if !got {
		t.Errorf("Failed to add Product %v", err)
	}

}
