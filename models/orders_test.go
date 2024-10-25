package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOrder(t *testing.T) {
	items := []OrderItem{
		OrderItem{ProductId: 2, Qty: 1, Price: 200.00},
		OrderItem{ProductId: 3, Qty: 2, Price: 250.00},
	}

	order := Order{
		CustomerId: 1,
		ShippingId: 1,
		Status:     0,
		Amount:     0,
		Items:      items,
	}

	orderId, err := AddOrderRecord(order)
	require.NoError(t, err)
	assert.Greater(t, orderId, 0)

	for _, item := range order.Items {
		item.OrderId = orderId
		ok, err := AddOrderItem(item)
		require.NoError(t, err)
		assert.Equal(t, ok, true)
	}

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

	assert.Equal(t, newOrder.Items[0].Price, 20000.00) // price in cents
	assert.Equal(t, newOrder.Amount, 70000.00)
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
