package models

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestAddOrder(t *testing.T) {

}

func TestMarshalOrder(t *testing.T) {
	items := []OrderItem{
		OrderItem{ProductId: 2, Qty: 1, Price: 200.00},
		OrderItem{ProductId: 3, Qty: 2, Price: 250.00},
	}

	order := Order{
		CustomerId: 1,
		Status:     0,
		Amount:     0,
		Items:      items,
	}

	res, _ := json.Marshal(order)

	var newOrder Order
	err := json.Unmarshal(res, &newOrder)
	if err != nil {
		t.Logf("ERR %v\n", err)
	}

	assert.Equal(t, newOrder.Amount, 70000)
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
