package models

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOrderRecord(t *testing.T) {
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

	_, orderId, err := AddOrder(order)
	require.NoError(t, err)
	assert.Greater(t, orderId, 0)

	orderItems, err := GetOrderItems(orderId)
	require.NoError(t, err)
	assert.Equal(t, 2, len(orderItems))
}

func TestMarshalOrder(t *testing.T) {
	items := []OrderItem{
		OrderItem{ProductId: 2, Qty: 1, Price: 200.00},
		OrderItem{ProductId: 3, Qty: 2, Price: 250.50},
	}

	order := Order{
		CustomerId: 1,
		Status:     0,
		Amount:     0,
		Items:      items,
	}

	res, err := json.Marshal(order)

	jsonStr := fmt.Sprintf("%s", res)
	expect := "{\"id\":0,\"shipping_id\":0,\"customer_id\":1,\"reference_code\":\"\",\"payment_reference\":\"\",\"amount_in_cents\":0,\"status\":0,\"orders\":[{\"id\":0,\"order_id\":0,\"product_id\":2,\"qty\":1,\"price\":200},{\"id\":0,\"order_id\":0,\"product_id\":3,\"qty\":2,\"price\":250.5}],\"amount\":701}"

	require.NoError(t, err)
	assert.Equal(t, expect, jsonStr)
	var newOrder Order
	err = json.Unmarshal(res, &newOrder)
	if err != nil {
		t.Logf("ERR %v\n", err)
	}

	assert.Equal(t, newOrder.Amount, 70100.00)
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
