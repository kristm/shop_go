package models

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestAddOrder(t *testing.T) {
	//Create customer
	//create order, pending status
	//crate order_products records
}

func TestUnmarshalOrderItem(t *testing.T) {
	data := []byte(`
		{
			"product_id": 2,
			"qty": 1,
			"price": 299.99
		}
	`)

	var orderItem OrderItem

	_ = json.Unmarshal(data, &orderItem)
	// price needs to be converted to cents before inserting to db
	assert.Equal(t, orderItem.Price, float64(29999))
}
